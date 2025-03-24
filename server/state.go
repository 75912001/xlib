package server

import (
	xconstants "github.com/75912001/xlib/constants"
	xcontrol "github.com/75912001/xlib/control"
	xlog "github.com/75912001/xlib/log"
	xtimer "github.com/75912001/xlib/timer"
	"runtime"
	"runtime/debug"
	"time"
)

func stateTimerPrint(timer xtimer.ITimer, l xlog.ILog) {
	defaultCallBack := xcontrol.NewCallBack(timeOut, timer, l)
	_ = timer.AddSecond(defaultCallBack, time.Now().Unix()+xconstants.ServerInfoTimeOutSec)
}

// 服务信息 打印
func timeOut(arg ...interface{}) error {
	s := debug.GCStats{}
	debug.ReadGCStats(&s)
	l := arg[1].(xlog.ILog)
	l.Infof("goroutineCnt:%v, numGC:%d, lastGC:%v, GCPauseTotal:%v",
		runtime.NumGoroutine(), s.NumGC, s.LastGC, s.PauseTotal)
	stateTimerPrint(arg[0].(xtimer.ITimer), l)
	return nil
}
