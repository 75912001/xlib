package server

import (
	xcontrol "github.com/75912001/xlib/control"
	xlog "github.com/75912001/xlib/log"
	xtimer "github.com/75912001/xlib/timer"
	"runtime"
	"runtime/debug"
	"sync/atomic"
	"time"
)

func stateTimerPrint(timer xtimer.ITimer, l xlog.ILog) {
	defaultCallBack := xcontrol.NewCallBack(timeOut, timer, l)
	_ = timer.AddSecond(defaultCallBack, time.Now().Unix()+ServerInfoTimeOutSec)
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

const StatusRunning = 0  // 服务状态：运行中
const StatusStopping = 1 // 服务状态：关闭中

var GQuitChan = make(chan bool)

var GServerStatus uint32

// IsServerStopping 服务是否关闭中
func IsServerStopping() bool {
	return atomic.LoadUint32(&GServerStatus) == StatusStopping
}

// SetServerStopping 设置为关闭中
func SetServerStopping() {
	atomic.StoreUint32(&GServerStatus, StatusStopping)
}
