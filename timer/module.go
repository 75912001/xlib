package timer

import (
	"context"
	xcontrol "github.com/75912001/xlib/control"
)

type ITimerSecond interface {
	AddSecond(callBackFunc xcontrol.ICallBack, expire int64) *Second
	DelSecond(second *Second)
}

type ITimerMillisecond interface {
	AddMillisecond(callBackFunc xcontrol.ICallBack, expireMillisecond int64) *millisecond
	DelMillisecond(millisecond *millisecond)
}

type ITimer interface {
	Start(ctx context.Context, opts ...*options) error
	Stop()
	ITimerSecond
	ITimerMillisecond
}

type EventTimerSecond struct {
	ISwitch   xcontrol.ISwitchButton
	ICallBack xcontrol.ICallBack
}

type EventTimerMillisecond struct {
	ISwitch   xcontrol.ISwitchButton
	ICallBack xcontrol.ICallBack
}
