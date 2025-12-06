package heartbeat

import (
	xcontrol "github.com/75912001/xlib/control"
	xtimer "github.com/75912001/xlib/timer"
	"time"
)

type HeartBeat struct {
	WaitID uint32         // 等待的id (开始为0, 收到第一次心跳,设置随机值,并返回给用户,用户下次使用该数值)
	Second *xtimer.Second // 定时器
	Timer  xtimer.ITimer
}

// 开始
//
//	callback.Parameters[0]: object 挂载对象
//	callback.Parameters[1]: timer xtimer.ITimer
//	callback.Parameters[2]: expire int64
func (p *HeartBeat) Start(callback xcontrol.ICallBack, out xcontrol.IOut) {
	parameters := callback.Get()
	p.Timer = parameters[1].(xtimer.ITimer)
	expire := parameters[2].(int64)
	p.Second = p.Timer.AddSecond(callback, time.Now().Unix()+expire, out)
}

func (p *HeartBeat) Stop() {
	if p.Second != nil {
		p.Timer.DelSecond(p.Second)
		p.Second = nil
	}
}
