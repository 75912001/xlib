package server

import (
	"context"
	"fmt"
	xconfig "github.com/75912001/xlib/config"
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xetcd "github.com/75912001/xlib/etcd"
	xlog "github.com/75912001/xlib/log"
	xcommon "github.com/75912001/xlib/net/common"
	xkcp "github.com/75912001/xlib/net/kcp"
	xtcp "github.com/75912001/xlib/net/tcp"
	xpprof "github.com/75912001/xlib/pprof"
	xruntime "github.com/75912001/xlib/runtime"
	xtime "github.com/75912001/xlib/time"
	xtimer "github.com/75912001/xlib/timer"
	xutil "github.com/75912001/xlib/util"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	BenchMgr xconfig.Mgr
	BenchSub *xconfig.Sub

	GroupID        uint32 // 组ID
	Name           string // 名称
	ID             uint32 // ID
	ExecutablePath string // 执行程序路径 // 程序所在路径(如为link,则为link所在的路径)

	Log     xlog.ILog
	TimeMgr *xtime.Mgr
	Timer   xtimer.ITimer

	Etcd    xetcd.IEtcd
	EtcdKey string // etcd key

	BusChannel          chan interface{} // 总线
	BusChannelWaitGroup sync.WaitGroup   // 总线等待

	QuitChan chan struct{} // 退出信号, 用于关闭服务

	TCPServer *xtcp.Server
	KCPServer *xkcp.Server
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

func (p *Server) PreStop() error {
	return xerror.NotImplemented
}

func (p *Server) Start(ctx context.Context,
	handler xcommon.IHandler,
	logCallbackFunc xlog.CallBackFunc,
	etcdCallbackFun xetcd.CallbackFun) (err error) {
	rand.Seed(time.Now().UnixNano())
	p.TimeMgr.Update()
	// 小端
	if !xutil.IsLittleEndian() {
		return errors.Errorf("system is bigEndian! %v", xruntime.Location())
	}
	// 开启UUID随机
	uuid.EnableRandPool()
	// 服务配置文件
	benchPath := path.Join(p.ExecutablePath, fmt.Sprintf("%v.%v.%v.%v",
		p.GroupID, p.Name, p.ID, xconfig.ServerConfigFileSuffix))
	content, err := os.ReadFile(benchPath)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	benchString := string(content)
	// 加载服务配置文件-root部分
	err = p.BenchMgr.Root.Parse(benchString)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	// 加载服务配置文件-公共部分
	err = p.BenchMgr.Config.Parse(benchString)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if false { // 从etcd获取配置项 todo menglc
		//client, err := clientv3.New(
		//	clientv3.Config{
		//		Endpoints:   p.BenchMgr.RootJson.Etcd.Addrs,
		//		DialTimeout: 5 * time.Second, // todo menglc 确定用途?
		//	},
		//)
		//if err != nil {
		//	return errors.WithMessage(err, xruntime.Location())
		//}
		//kv := clientv3.NewKV(client)
		//key := fmt.Sprintf("/%v/%v/%v/%v/%v",
		//	*p.BenchMgr.Json.Base.ProjectName, xetcd.WatchMsgTypeServiceBench, p.GroupID, p.Name, p.ID)
		//getResponse, err := kv.Get(ctx, key, clientv3.WithPrefix())
		//if err != nil {
		//	return errors.WithMessage(err, xruntime.Location())
		//}
		//if len(getResponse.Kvs) != 1 {
		//	return errors.WithMessagef(xerror.Config, "%v %v %v", key, getResponse.Kvs, xruntime.Location())
		//}
		//benchString = string(getResponse.Kvs[0].Value)
		//xlog.PrintfInfo(benchString)
	}
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	switch *p.BenchMgr.Config.Base.RunMode {
	case 0:
		xruntime.SetRunMode(xruntime.RunModeRelease)
	case 1:
		xruntime.SetRunMode(xruntime.RunModeDebug)
	default:
		return errors.Errorf("runMode err:%v %v", *p.BenchMgr.Config.Base.RunMode, xruntime.Location())
	}
	// GoMaxProcess
	previous := runtime.GOMAXPROCS(*p.BenchMgr.Config.Base.GoMaxProcess)
	xlog.PrintfInfo("go max process new:%v, previous setting:%v",
		*p.BenchMgr.Config.Base.GoMaxProcess, previous)
	// 日志
	p.Log, err = xlog.NewMgr(xlog.NewOptions().
		WithLevel(*p.BenchMgr.Config.Base.LogLevel).
		WithAbsPath(*p.BenchMgr.Config.Base.LogAbsPath).
		WithNamePrefix(fmt.Sprintf("%v.%v.%v", p.GroupID, p.Name, p.ID)).
		WithLevelCallBack(logCallbackFunc, xlog.LevelFatal, xlog.LevelError, xlog.LevelWarn),
	)
	if err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	// 加载服务配置文件-子项部分
	if p.BenchSub != nil {
		err = p.BenchSub.Unmarshal(benchString)
		if err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
	}
	// eventChan
	p.BusChannel = make(chan interface{}, *p.BenchMgr.Config.Base.BusChannelCapacity)
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
	if p.BenchMgr.Config.Base.PprofHttpPort != nil {
		xpprof.StartHTTPprof(fmt.Sprintf("0.0.0.0:%d", *p.BenchMgr.Config.Base.PprofHttpPort))
	}
	// 全局定时器
	if p.BenchMgr.Config.Timer.ScanSecondDuration != nil || p.BenchMgr.Config.Timer.ScanMillisecondDuration != nil {
		p.Timer = xtimer.NewTimer()
		err = p.Timer.Start(ctx,
			xtimer.NewOptions().
				WithScanSecondDuration(*p.BenchMgr.Config.Timer.ScanSecondDuration).
				WithScanMillisecondDuration(*p.BenchMgr.Config.Timer.ScanMillisecondDuration).
				WithOutgoingTimerOutChan(p.BusChannel),
		)
		if err != nil {
			return errors.Errorf("timer Start err:%v %v", err, xruntime.Location())
		}
	}
	// etcd
	p.EtcdKey = xetcd.GenKey(*p.BenchMgr.Config.Base.ProjectName, xetcd.WatchMsgTypeServer, p.GroupID, p.Name, p.ID)
	defaultEtcd := xetcd.NewEtcd(
		xetcd.NewOptions().
			WithAddrs(p.BenchMgr.Root.Etcd.Addrs).
			WithTTL(*p.BenchMgr.Root.Etcd.TTL).
			WithWatchKeyPrefix(xetcd.GenPrefixKey(*p.BenchMgr.Config.Base.ProjectName)).
			WithKey(p.EtcdKey).
			WithValue(
				&xetcd.ValueJson{
					ServerNet:     p.BenchMgr.Config.ServerNet,
					Version:       *p.BenchMgr.Config.Base.Version,
					AvailableLoad: *p.BenchMgr.Config.Base.AvailableLoad,
					SecondOffset:  0,
				},
			).
			WithEventChan(p.BusChannel),
	)
	defaultEtcd.CallbackFun = etcdCallbackFun
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
	p.Timer.AddSecond(xcontrol.NewCallBack(etcdReportFunction, p), p.TimeMgr.ShadowTimestamp()+xetcd.ReportIntervalSecondDefault)
	// 网络服务
	for _, element := range p.BenchMgr.Config.ServerNet {
		if len(*element.Addr) != 0 {
			switch *element.Type {
			case "tcp": // 启动 TCP 服务
				p.TCPServer = xtcp.NewServer(handler)
				if err = p.TCPServer.Start(ctx,
					xtcp.NewServerOptions().
						WithListenAddress(*element.Addr).
						WithEventChan(p.BusChannel).
						WithSendChanCapacity(*p.BenchMgr.Config.Base.SendChannelCapacity),
				); err != nil {
					return errors.WithMessage(err, xruntime.Location())
				}
			case "kcp":
				p.KCPServer = nil // todo menglc xkcp.NewServer(handler)
				return errors.WithMessage(xerror.NotImplemented, xruntime.Location())
			default:
				return errors.WithMessage(xerror.NotImplemented, xruntime.Location())
			}
		}
	}

	stateTimerPrint(p.Timer, p.Log)
	return nil
}

func (p *Server) Stop() (err error) {
	if p.TCPServer != nil {
		p.TCPServer.Stop()
	}
	if p.KCPServer != nil {
		p.KCPServer.Stop()
	}
	return nil
}
