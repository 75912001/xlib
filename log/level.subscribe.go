package log

import (
	xcontrol "github.com/75912001/xlib/control"
	xmap "github.com/75912001/xlib/map"
)

type levelSubscribe struct {
	subMapMgr *xmap.MapMgr[uint32, struct{}]
	callBack  xcontrol.ICallBack
}

func newLevelSubscribe() *levelSubscribe {
	return &levelSubscribe{
		subMapMgr: xmap.NewMapMgr[uint32, struct{}](),
	}
}

// 是否订阅
func (p *levelSubscribe) isSubscribe(level uint32) bool {
	return p.subMapMgr.IsExist(level)
}
