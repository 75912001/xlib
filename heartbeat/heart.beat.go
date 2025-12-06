package heartbeat

import (
	xcontrol "github.com/75912001/xlib/control"
)

// IHeartBeat 心跳接口
type IHeartBeat interface {
	Start(callback xcontrol.ICallBack, out xcontrol.IOut) // 开启
	Stop()                                                // 停止
	Timeout(args ...any) error
}
