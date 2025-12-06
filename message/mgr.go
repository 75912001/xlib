package message

import (
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xmap "github.com/75912001/xlib/map"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

// Mgr 管理器
type Mgr struct {
	mapMgr *xmap.MapMgr[uint32, IMessage] // key: messageID, value: IMessage
}

func NewMgr() *Mgr {
	return &Mgr{
		mapMgr: xmap.NewMapMgr[uint32, IMessage](),
	}
}

// Register 注册消息
//
//	重复会 panic
func (p *Mgr) Register(messageID uint32, opts ...*options) {
	if pb := p.Find(messageID); pb != nil {
		xlog.PrintErr(xerror.Exist, "%v messageID:%#x %v", xruntime.Location(), messageID, messageID)
		panic(errors.WithMessagef(xerror.Exist, "%v messageID:%#x %v",
			xruntime.Location(), messageID, messageID))
	}
	opt := merge(opts...)
	if err := configure(opt); err != nil {
		xlog.PrintErr(xerror.Param, "%v messageID:%#x %v", xruntime.Location(), messageID, messageID)
		panic(errors.WithMessagef(xerror.Param, "%v messageID:%#x %v",
			xruntime.Location(), messageID, messageID))
	}
	p.mapMgr.Add(messageID, newDefaultMessage(opt))
}

func (p *Mgr) Find(messageID uint32) IMessage {
	msg, ok := p.mapMgr.Find(messageID)
	if !ok {
		return nil
	}
	return msg
}

// Replace 替换/覆盖(Override)
//
//	配置错误会 panic
func (p *Mgr) Replace(messageID uint32, opt *options) {
	if err := configure(opt); err != nil {
		xlog.PrintErr(xerror.Param, "%v messageID:%#x %v", xruntime.Location(), messageID, messageID)
		panic(errors.WithMessagef(xerror.Param, "%v messageID:%#x %v",
			xruntime.Location(), messageID, messageID))
	}
	p.mapMgr.Add(messageID, newDefaultMessage(opt))
}
