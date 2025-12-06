package timer

import (
	"context"
	xcontrol "github.com/75912001/xlib/control"
)

var GTimer ITimer

type ITimerSecond interface {
	AddSecond(callBackFunc xcontrol.ICallBack, expire int64, out xcontrol.IOut) *Second // callBackFunc: 到期-回调函数 expire: 到期 时间戳(秒) out: 到期-输出
	DelSecond(second *Second)
}

type ITimerMillisecond interface {
	AddMillisecond(callBackFunc xcontrol.ICallBack, expireMillisecond int64, out xcontrol.IOut) *Millisecond // callBackFunc: 到期-回调函数 expireMillisecond: 到期 时间戳(毫秒) out: 到期-输出
	DelMillisecond(millisecond *Millisecond)
}

type ITimer interface {
	Start(ctx context.Context) error
	Stop()
	ITimerSecond
	ITimerMillisecond
}
