package server

import (
	"context"
	"crypto/sha1"
	"fmt"
	xconfig "github.com/75912001/xlib/config"
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xetcd "github.com/75912001/xlib/etcd"
	xetcdconstants "github.com/75912001/xlib/etcd/constants"
	xgrpcprotoregistry "github.com/75912001/xlib/grpc/proto/registry"
	xgrpcselector "github.com/75912001/xlib/grpc/selector"
	xgrpc "github.com/75912001/xlib/grpc/server"
	xlog "github.com/75912001/xlib/log"
	xnetcommon "github.com/75912001/xlib/net/common"
	xnetkcp "github.com/75912001/xlib/net/kcp"
	xnettcp "github.com/75912001/xlib/net/tcp"
	xnetwebsocket "github.com/75912001/xlib/net/websocket"
	xpprof "github.com/75912001/xlib/pprof"
	xruntime "github.com/75912001/xlib/runtime"
	xruntimeconstants "github.com/75912001/xlib/runtime/constants"
	xserverresources "github.com/75912001/xlib/server/resources"
	xtimer "github.com/75912001/xlib/timer"
	"github.com/pkg/errors"
	"github.com/xdg-go/pbkdf2"
	"github.com/xtaci/kcp-go/v5"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

type Server struct {
	actor *Actor

	QuitChan chan struct{} // 退出信号, 用于关闭服务

	TCPServer  *xnettcp.Server
	KCPServer  *xnetkcp.Server
	GRPCServer *xgrpc.Server
	WebSocket  *xnetwebsocket.Server

	Options *Options
	Derived IServer // 服务实例
}

// NewServer 创建服务
// args: [0:程序名称] [1:配置文件绝对路径]
func NewServer(args []string) *Server {
	s := &Server{
		QuitChan: make(chan struct{}),
	}
	// 程序所在路径(如为link,则为link所在的路径)
	argNum := len(args)
	const neededArgsNumber = 2
	if argNum != neededArgsNumber {
		xlog.PrintfErr("the number of parameters is incorrect, needed %v, but %v.", neededArgsNumber, argNum)
		return nil
	}
	// 配置文件路径
	configPath := args[1]
	// 加载服务配置文件
	err := xconfig.GConfigMgr.Parse(configPath)
	if err != nil {
		xlog.PrintfErr("parse config file failed,configPath %v %v", configPath, err)
		return nil
	}
	s.actor = NewActor(uint64(*xconfig.GConfigMgr.Base.ServerID), s.behavior)
	return s
}

func (p *Server) GetOptions() (opt *Options) {
	return p.Options
}

func (p *Server) GetActor() *Actor {
	return p.actor
}

func (p *Server) PreStart(ctx context.Context, opts ...*Options) error {
	p.Options = mergeOptions(opts...)
	if err := configure(p.Options); err != nil {
		return errors.WithMessagef(err, "configure err. %v", xruntime.Location())
	}
	switch *xconfig.GConfigMgr.Base.RunMode {
	case 0:
		xruntime.SetRunMode(xruntimeconstants.RunModeRelease)
	case 1:
		xruntime.SetRunMode(xruntimeconstants.RunModeDebug)
	default:
		return errors.WithMessagef(xerror.Param, "runMode err. runMode:%v %v", *xconfig.GConfigMgr.Base.RunMode, xruntime.Location())
	}

	xserverresources.GResources.SetAvailableLoad(*xconfig.GConfigMgr.Base.AvailableLoad)

	// GoMaxProcess
	previous := runtime.GOMAXPROCS(*xconfig.GConfigMgr.Base.GoMaxProcess)
	xlog.PrintfInfo("go max process new:%v, previous setting:%v",
		*xconfig.GConfigMgr.Base.GoMaxProcess, previous)

	var err error
	// 日志
	xlog.GLog, err = xlog.NewMgr(xlog.NewOptions().
		WithLevel(*xconfig.GConfigMgr.Log.Level).
		WithAbsPath(*xconfig.GConfigMgr.Log.AbsPath).
		WithNamePrefix(fmt.Sprintf("%v.%v.%v", *xconfig.GConfigMgr.Base.GroupID, *xconfig.GConfigMgr.Base.Name, *xconfig.GConfigMgr.Base.ServerID)).
		WithLevelCallBack(p.Options.LogCallback, xlog.LevelFatal, xlog.LevelError, xlog.LevelWarn, xlog.LevelInfo, xlog.LevelDebug, xlog.LevelTrace),
	)
	if err != nil {
		return errors.WithMessagef(err, "log newMgr err. %v", xruntime.Location())
	}
	xlog.GLog.Info("========== server start - log ==========")

	// 初始化 proto 扩展配置
	xgrpcprotoregistry.Init()
	xgrpcselector.Init()
	// grpc 服务
	if xconfig.GConfigMgr.Grpc.IsEnabled() {
		p.GRPCServer = xgrpc.NewServer()
	}

	p.actor.Start()

	// 全局定时器
	{
		xtimer.GTimer = xtimer.NewTimer()
		err = xtimer.GTimer.Start(ctx)
		if err != nil {
			return errors.Errorf("timer Start err:%v %v", err, xruntime.Location())
		}
	}

	// 是否开启http采集分析
	if xconfig.GConfigMgr.Base.PprofHttpPort != nil {
		xpprof.StartHTTPprof(fmt.Sprintf("0.0.0.0:%d", *xconfig.GConfigMgr.Base.PprofHttpPort))
	}

	return nil
}

func (p *Server) Start(ctx context.Context) (err error) {
	////////////////////////////////////////////////////////////
	// grpc 服务
	if xconfig.GConfigMgr.Grpc.IsEnabled() {
		err = p.GRPCServer.Start(*xconfig.GConfigMgr.Grpc.ListenAddr)
		if err != nil {
			return errors.WithMessagef(err, "grpc server start err. %v", xruntime.Location())
		}
	}
	////////////////////////////////////////////////////////////
	// 网络服务
	for _, element := range xconfig.GConfigMgr.Net {
		switch *element.Type {
		case xnetcommon.ServerNetTypeNameTCP: // 启动 TCP 服务
			p.TCPServer = xnettcp.NewServer(p.Options.TCPHandler)
			serverOptions := xnettcp.NewServerOptions().
				WithListenAddress(*element.ListenAddr).
				WithIOut(p.GetActor()).
				WithSendChanCapacity(*xconfig.GConfigMgr.Base.SendChannelCapacity).
				WithHeaderStrategy(p.Options.HeaderStrategy)
			serverOptions.WithNewPacketLimitFunc(xnetcommon.NewPackLimitDefault).
				WithMaxCntPerSec(*xconfig.GConfigMgr.Base.PacketLimitRecvCntPreSecond)
			if err = p.TCPServer.Start(ctx, serverOptions); err != nil {
				return errors.WithMessagef(err, "tcp server start err. %v", xruntime.Location())
			}
		case xnetcommon.ServerNetTypeNameKCP:
			p.KCPServer = xnetkcp.NewServer(p.Options.KCPHandler)
			var blockCrypt kcp.BlockCrypt
			key := pbkdf2.Key([]byte(*xconfig.GConfigMgr.KCP.Password), []byte(*xconfig.GConfigMgr.KCP.Salt), 1024, 32, sha1.New)
			blockCrypt, err = kcp.NewAESBlockCrypt(key)
			if err != nil {
				return errors.WithMessagef(err, "kcp server start err. %v", xruntime.Location())
			}
			kcpOpts := xnetkcp.NewOptions()
			kcpOpts.WithListenAddress(*element.ListenAddr).
				WithIOut(p.GetActor()).
				WithSendChanCapacity(*xconfig.GConfigMgr.Base.SendChannelCapacity).
				WithHeaderStrategy(p.Options.HeaderStrategy).
				WithNewPacketLimitFunc(xnetcommon.NewPackLimitDefault)
			kcpOpts.WithNewPacketLimitFunc(xnetcommon.NewPackLimitDefault).
				WithMaxCntPerSec(*xconfig.GConfigMgr.Base.PacketLimitRecvCntPreSecond)
			kcpOpts.WithBlockCrypt(blockCrypt).
				WithFEC(true)
			if err = p.KCPServer.Start(ctx, kcpOpts); err != nil {
				return errors.WithMessagef(err, "kcp server start err. %v", xruntime.Location())
			}
		case xnetcommon.ServerNetTypeNameWebSocket:
			p.WebSocket = xnetwebsocket.NewServer(p.Options.WebsocketHandler)
			serverOptions := xnetwebsocket.NewServerOptions().
				WithPattern(*element.Pattern).
				WithListenAddress(*element.ListenAddr).
				WithIOut(p.GetActor()).
				WithSendChanCapacity(*xconfig.GConfigMgr.Base.SendChannelCapacity)
			serverOptions.WithNewPacketLimitFunc(xnetcommon.NewPackLimitDefault).
				WithMaxCntPerSec(*xconfig.GConfigMgr.Base.PacketLimitRecvCntPreSecond)
			if err = p.WebSocket.Start(ctx, serverOptions); err != nil {
				return errors.WithMessagef(err, "websocket server start err. %v", xruntime.Location())
			}
		default:
			return errors.WithMessagef(xerror.NotImplemented, "server net type not implemented. %v", xruntime.Location())
		}
	}

	stateTimerPrint(xtimer.GTimer, xlog.GLog, p.GetActor())
	////////////////////////////////////////////////////////////
	// etcd
	etcdKey := xetcd.GenKey(*xconfig.GConfigMgr.Base.ProjectName,
		xetcdconstants.WatchMsgTypeServer,
		*xconfig.GConfigMgr.Base.GroupID, *xconfig.GConfigMgr.Base.Name, *xconfig.GConfigMgr.Base.ServerID)
	value := p.genEtcdValue()

	opt := xetcd.MergeOptions(p.Options.Etcd)
	defaultEtcd := xetcd.NewEtcd(
		xetcd.NewOptions().
			WithEndpoints(xconfig.GConfigMgr.Etcd.Endpoints).
			WithTTL(*xconfig.GConfigMgr.Etcd.TTL).
			WithWatchKeyPrefix(xetcd.GenPrefixKey(*xconfig.GConfigMgr.Base.ProjectName)).
			WithKey(etcdKey).
			WithIOut(p.GetActor()),
		opt,
	)
	xetcd.GEtcd = defaultEtcd
	if err = xetcd.GEtcd.Start(ctx, value); err != nil {
		return errors.WithMessagef(err, "etcd start err. %v", xruntime.Location())
	}
	// 续租
	err = xetcd.GEtcd.KeepAlive(ctx)
	if err != nil {
		return errors.WithMessagef(err, "etcd keepAlive err. %v", xruntime.Location())
	}
	// etcd-定时上报
	xtimer.GTimer.AddSecond(xcontrol.NewCallBack(etcdReportFunction, p),
		time.Now().Unix()+ReportIntervalSecondDefault,
		p.GetActor(),
	)
	////////////////////////////////////////////////////////////
	runtime.GC()
	return nil
}

func (p *Server) genEtcdValue() string {
	valueJson := &xetcd.ValueJson{
		Version:       *xconfig.GConfigMgr.Base.Version,
		AvailableLoad: xserverresources.GResources.GetAvailableLoad(),
		SecondOffset:  0,
	}
	for _, v := range xconfig.GConfigMgr.Net {
		valueJson.ServerNet = append(valueJson.ServerNet,
			&xetcd.ServerNet{
				Addr: v.ExternalAddr,
				Name: v.Name,
				Type: v.Type,
			},
		)
	}
	if xconfig.GConfigMgr.Grpc.IsEnabled() {
		valueJson.GrpcService = &xetcd.GrpcService{
			PackageName: xconfig.GConfigMgr.Grpc.PackageName,
			ServiceName: xconfig.GConfigMgr.Grpc.ServiceName,
			Addr:        xconfig.GConfigMgr.Grpc.ExternalAddr,
		}
	}
	return xetcd.ValueJson2String(valueJson)
}

func (p *Server) PostStart() error {
	// 退出服务
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	select {
	case <-p.QuitChan:
		xlog.GLog.Warn("Server will shutdown in a few seconds")
	case s := <-sigChan:
		xlog.GLog.Warnf("Server got signal: %s, shutting down...", s)
	}
	err := p.Derived.PreStop()
	if err != nil {
		xlog.GLog.Warn("pre stop err:%v ", err)
	}
	// 设置为关闭中
	SetServerStopping()

	p.actor.Stop()

	err = p.Derived.Stop()
	if err != nil {
		xlog.GLog.Warn("server stop err:%v ", err)
		return errors.WithMessagef(err, "server stop err. %v", xruntime.Location())
	}
	return nil
}

func (p *Server) PreStop() error {
	return xerror.NotImplemented
}

func (p *Server) Stop() (err error) {
	err = xetcd.GEtcd.Stop()
	if err != nil {
		xlog.GLog.Errorf("etcd stop err:%v", err)
	}
	if xconfig.GConfigMgr.Grpc.IsEnabled() {
		err = p.GRPCServer.Stop()
		if err != nil {
			xlog.GLog.Errorf("grpc server stop err:%v", err)
		}
	}
	if p.TCPServer != nil {
		p.TCPServer.Stop()
	}
	if p.KCPServer != nil {
		p.KCPServer.Stop()
	}
	if p.WebSocket != nil {
		p.WebSocket.Stop()
	}
	if p.GRPCServer != nil {
		_ = p.GRPCServer.Stop()
	}

	xtimer.GTimer.Stop()
	return nil
}
