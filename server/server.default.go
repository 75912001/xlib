package server

import (
	"context"
	"crypto/sha1"
	"fmt"
	xconfig "github.com/75912001/xlib/config"
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xetcd "github.com/75912001/xlib/etcd"
	xlog "github.com/75912001/xlib/log"
	xnetcommon "github.com/75912001/xlib/net/common"
	xnetkcp "github.com/75912001/xlib/net/kcp"
	xnettcp "github.com/75912001/xlib/net/tcp"
	xpprof "github.com/75912001/xlib/pprof"
	xruntime "github.com/75912001/xlib/runtime"
	xtime "github.com/75912001/xlib/time"
	xtimer "github.com/75912001/xlib/timer"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/xdg-go/pbkdf2"
	"github.com/xtaci/kcp-go/v5"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	ConfigMgr xconfig.Mgr
	ConfigSub *xconfig.Sub

	GroupID        uint32 // 组ID
	Name           string // 名称
	ID             uint32 // ID
	ExecutablePath string // 执行程序路径 // 程序所在路径(如为link,则为link所在的路径)
	AvailableLoad  uint32 // 可用负载

	Log     xlog.ILog
	TimeMgr *xtime.Mgr
	Timer   xtimer.ITimer

	Etcd    xetcd.IEtcd
	EtcdKey string // etcd key

	BusChannel          chan interface{} // 总线
	BusChannelWaitGroup sync.WaitGroup   // 总线等待

	QuitChan chan struct{} // 退出信号, 用于关闭服务

	TCPServer *xnettcp.Server
	KCPServer *xnetkcp.Server

	options      *ServerOptions
	ServerObject IServer // 服务实例
}

// NewServer 创建服务
// args: [进程名称, 组ID, 服务名, 服务ID]
func NewServer(args []string) *Server {
	s := &Server{
		TimeMgr:  xtime.NewMgr(),
		QuitChan: make(chan struct{}),
	}
	// 程序所在路径(如为link,则为link所在的路径)
	if executablePath, err := xruntime.GetExecutablePath(); err != nil {
		xlog.PrintErr(err, xruntime.Location())
		return nil
	} else {
		s.ExecutablePath = executablePath
	}
	argNum := len(args)
	const neededArgsNumber = 4
	if argNum != neededArgsNumber {
		xlog.PrintfErr("the number of parameters is incorrect, needed %v, but %v.", neededArgsNumber, argNum)
		return nil
	}
	{ // 解析启动参数
		groupID, err := strconv.ParseUint(args[1], 10, 32)
		if err != nil {
			xlog.PrintErr("groupID err:", err)
			return nil
		}
		s.GroupID = uint32(groupID)
		s.Name = args[2]
		serviceID, err := strconv.ParseUint(args[3], 10, 32)
		if err != nil {
			xlog.PrintErr("serviceID err", err)
			return nil
		}
		s.ID = uint32(serviceID)
		xlog.PrintfInfo("groupID:%v name:%v, serviceID:%v",
			s.GroupID, s.Name, s.ID)
	}
	return s
}

func (p *Server) Start(ctx context.Context, opts ...*ServerOptions) (err error) {
	p.options = mergeServerOptions(opts...)
	if err := serverConfigure(p.options); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	rand.Seed(time.Now().UnixNano())
	p.TimeMgr.Update()
	// 开启UUID随机
	uuid.EnableRandPool()
	// 服务配置文件
	configPath := path.Join(p.ExecutablePath, fmt.Sprintf("%v.%v.%v.%v",
		p.GroupID, p.Name, p.ID, xconfig.ServerConfigFileSuffix))
	content, err := os.ReadFile(configPath)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	configString := string(content)
	// 加载服务配置文件-root部分
	err = p.ConfigMgr.Root.Parse(configString)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	// 加载服务配置文件-公共部分
	err = p.ConfigMgr.Config.Parse(configString)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if false { // 从etcd获取配置项 todo menglc
		//client, err := clientv3.New(
		//	clientv3.Config{
		//		Endpoints:   p.ConfigMgr.RootJson.Etcd.Addrs,
		//		DialTimeout: 5 * time.Second, // todo menglc 确定用途?
		//	},
		//)
		//if err != nil {
		//	return errors.WithMessage(err, xruntime.Location())
		//}
		//kv := clientv3.NewKV(client)
		//key := fmt.Sprintf("/%v/%v/%v/%v/%v",
		//	*p.ConfigMgr.Json.Base.ProjectName, xetcd.WatchMsgTypeServiceBench, p.GroupID, p.Name, p.ID)
		//getResponse, err := kv.Get(ctx, key, clientv3.WithPrefix())
		//if err != nil {
		//	return errors.WithMessage(err, xruntime.Location())
		//}
		//if len(getResponse.Kvs) != 1 {
		//	return errors.WithMessagef(xerror.Config, "%v %v %v", key, getResponse.Kvs, xruntime.Location())
		//}
		//configString = string(getResponse.Kvs[0].Value)
		//xlog.PrintfInfo(configString)
	}
	switch *p.ConfigMgr.Config.Base.RunMode {
	case 0:
		xruntime.SetRunMode(xruntime.RunModeRelease)
	case 1:
		xruntime.SetRunMode(xruntime.RunModeDebug)
	default:
		return errors.Errorf("runMode err:%v %v", *p.ConfigMgr.Config.Base.RunMode, xruntime.Location())
	}
	p.AvailableLoad = *p.ConfigMgr.Config.Base.AvailableLoad
	// GoMaxProcess
	previous := runtime.GOMAXPROCS(*p.ConfigMgr.Config.Base.GoMaxProcess)
	xlog.PrintfInfo("go max process new:%v, previous setting:%v",
		*p.ConfigMgr.Config.Base.GoMaxProcess, previous)
	// 日志
	p.Log, err = xlog.NewMgr(xlog.NewOptions().
		WithLevel(*p.ConfigMgr.Config.Base.LogLevel).
		WithAbsPath(*p.ConfigMgr.Config.Base.LogAbsPath).
		WithNamePrefix(fmt.Sprintf("%v.%v.%v", p.GroupID, p.Name, p.ID)).
		WithLevelCallBack(p.options.LogCallbackFunc, xlog.LevelFatal, xlog.LevelError, xlog.LevelWarn),
	)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	// 加载服务配置文件-子项部分
	if p.ConfigSub != nil {
		err = p.ConfigSub.Unmarshal(configString)
		if err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
	}
	// eventChan
	p.BusChannel = make(chan interface{}, *p.ConfigMgr.Config.Base.BusChannelCapacity)
	go func() {
		defer func() {
			p.BusChannelWaitGroup.Done()
			// 主事件 channel 报错 不 recover
			p.Log.Infof(xerror.GoroutineDone.Error())
		}()
		p.BusChannelWaitGroup.Add(1)
		_ = p.Handle()
	}()
	// 是否开启http采集分析
	if p.ConfigMgr.Config.Base.PprofHttpPort != nil {
		xpprof.StartHTTPprof(fmt.Sprintf("0.0.0.0:%d", *p.ConfigMgr.Config.Base.PprofHttpPort))
	}
	// 全局定时器
	if p.ConfigMgr.Config.Timer.ScanSecondDuration != nil || p.ConfigMgr.Config.Timer.ScanMillisecondDuration != nil {
		p.Timer = xtimer.NewTimer()
		err = p.Timer.Start(ctx,
			xtimer.NewOptions().
				WithScanSecondDuration(*p.ConfigMgr.Config.Timer.ScanSecondDuration).
				WithScanMillisecondDuration(*p.ConfigMgr.Config.Timer.ScanMillisecondDuration).
				WithOutgoingTimerOutChan(p.BusChannel),
		)
		if err != nil {
			return errors.Errorf("timer Start err:%v %v", err, xruntime.Location())
		}
	}
	// etcd
	p.EtcdKey = xetcd.GenKey(*p.ConfigMgr.Config.Base.ProjectName, xetcd.WatchMsgTypeServer, p.GroupID, p.Name, p.ID)
	defaultEtcd := xetcd.NewEtcd(
		xetcd.NewOptions().
			WithAddrs(p.ConfigMgr.Root.Etcd.Addrs).
			WithTTL(*p.ConfigMgr.Root.Etcd.TTL).
			WithWatchKeyPrefix(xetcd.GenPrefixKey(*p.ConfigMgr.Config.Base.ProjectName)).
			WithKey(p.EtcdKey).
			WithValue(
				&xetcd.ValueJson{
					ServerNet:     p.ConfigMgr.Config.ServerNet,
					Version:       *p.ConfigMgr.Config.Base.Version,
					AvailableLoad: p.AvailableLoad,
					SecondOffset:  0,
				},
			).
			WithEventChan(p.BusChannel),
	)
	defaultEtcd.CallbackFun = p.options.ETCDCallbackFun
	p.Etcd = defaultEtcd
	if err = p.Etcd.Start(ctx); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	// 续租
	err = defaultEtcd.KeepAlive(ctx)
	if err != nil {
		return errors.WithMessagef(err, xruntime.Location())
	}
	// etcd-定时上报
	p.Timer.AddSecond(xcontrol.NewCallBack(etcdReportFunction, p), p.TimeMgr.ShadowTimestamp()+ReportIntervalSecondDefault)
	// 网络服务
	for _, element := range p.ConfigMgr.Config.ServerNet {
		if len(*element.Addr) != 0 {
			switch *element.Type {
			case xnetcommon.ServerNetTypeNameTCP: // 启动 TCP 服务
				p.TCPServer = xnettcp.NewServer(p.options.TCPHandler)
				if err = p.TCPServer.Start(ctx,
					xnettcp.NewServerOptions().
						WithListenAddress(*element.Addr).
						WithEventChan(p.BusChannel).
						WithSendChanCapacity(*p.ConfigMgr.Config.Base.SendChannelCapacity),
				); err != nil {
					return errors.WithMessage(err, xruntime.Location())
				}
			case xnetcommon.ServerNetTypeNameKCP:
				p.KCPServer = xnetkcp.NewServer(p.options.KCPHandler)
				var blockCrypt kcp.BlockCrypt
				if true { // 使用默认加密方式
					key := pbkdf2.Key([]byte("demo.pass"), []byte("demo.salt"), 1024, 32, sha1.New)
					blockCrypt, err = kcp.NewAESBlockCrypt(key)
					if err != nil {
						return errors.WithMessage(err, xruntime.Location())
					}
				}
				if err = p.KCPServer.Start(ctx,
					xnetkcp.NewOptions().
						WithListenAddress(*element.Addr).
						WithEventChan(p.BusChannel).
						WithSendChanCapacity(*p.ConfigMgr.Config.Base.SendChannelCapacity).
						WithWriteBuffer(1024*1024).
						WithReadBuffer(1024*1024).
						WithBlockCrypt(blockCrypt).
						WithMTUBytes(1350).
						WithWindowSize(512).
						WithFEC(true).
						WithAckNoDelay(true),
				); err != nil {
					return errors.WithMessage(err, xruntime.Location())
				}
			default:
				return errors.WithMessage(xerror.NotImplemented, xruntime.Location())
			}
		}
	}

	stateTimerPrint(p.Timer, p.Log)

	runtime.GC()

	// 退出服务
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	select {
	case <-p.QuitChan:
		p.Log.Warn("Server will shutdown in a few seconds")
	case s := <-sigChan:
		p.Log.Warnf("Server got signal: %s, shutting down...", s)
	}
	err = p.ServerObject.PreStop()
	if err != nil {
		p.Log.Warn("pre stop err:%v ", err)
	}
	// 设置为关闭中
	SetServerStopping()
	// 定时检查事件总线是否消费完成
	go p.checkGBusChannel()
	// 等待GEventChan处理结束
	p.BusChannelWaitGroup.Wait()
	err = p.ServerObject.Stop()
	if err != nil {
		p.Log.Warn("server stop err:%v ", err)
	}
	return nil
}

func (p *Server) PreStop() error {
	return xerror.NotImplemented
}

func (p *Server) Stop() (err error) {
	err = p.Etcd.Stop()
	if err != nil {
		p.Log.Errorf("etcd stop err:%v", err)
	}
	if p.TCPServer != nil {
		p.TCPServer.Stop()
	}
	if p.KCPServer != nil {
		p.KCPServer.Stop()
	}
	p.Timer.Stop()
	return nil
}

func (p *Server) checkGBusChannel() {
	p.Log.Warn("start checkGBusChannel timer")

	idleDuration := 500 * time.Millisecond
	idleDelay := time.NewTimer(idleDuration)
	defer func() {
		idleDelay.Stop()
	}()

	for {
		select {
		case <-idleDelay.C:
			idleDelay.Reset(idleDuration)
			GBusChannelQuitCheck <- struct{}{}
			p.Log.Warn("send to GBusChannelQuitCheck")
		}
	}
}
