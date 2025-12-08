package server

import (
	xconfig "github.com/75912001/xlib/config"
	xcontrol "github.com/75912001/xlib/control"
	xlog "github.com/75912001/xlib/log"
	xserverconstants "github.com/75912001/xlib/server/constants"
	xserverresources "github.com/75912001/xlib/server/resources"
	xtimer "github.com/75912001/xlib/timer"
	"runtime"
	"runtime/debug"
	"sync/atomic"
	"time"
)

func stateTimerPrint(timer xtimer.ITimer, l xlog.ILog, out xcontrol.IOut) {
	defaultCallBack := xcontrol.NewCallBack(timeOut, timer, l, out)
	_ = timer.AddSecond(defaultCallBack, time.Now().Unix()+xserverconstants.ServerInfoTimeOutSec, out)
}

// 服务信息 打印
func timeOut(args ...any) error {
	s := debug.GCStats{}
	debug.ReadGCStats(&s)
	l := args[1].(xlog.ILog)
	l.Infof("goroutineCnt:%v, numGC:%d, lastGC:%v, GCPauseTotal:%v availableLoad/AvailableLoad:%v/%v",
		runtime.NumGoroutine(), s.NumGC, s.LastGC, s.PauseTotal, xserverresources.GResources.GetAvailableLoad(),
		*xconfig.GConfigMgr.Base.AvailableLoad)
	stateTimerPrint(args[0].(xtimer.ITimer), l, args[2].(xcontrol.IOut))
	return nil
}

var GServerStatus uint32

// IsServerStopping 服务是否关闭中
func IsServerStopping() bool {
	return atomic.LoadUint32(&GServerStatus) == xserverconstants.StatusStopping
}

// SetServerStopping 设置为关闭中
func SetServerStopping() {
	atomic.StoreUint32(&GServerStatus, xserverconstants.StatusStopping)
}
