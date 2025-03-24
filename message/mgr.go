package message

import (
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

// Mgr 管理器
type Mgr struct {
	messageMap map[uint32]IMessage
}

// Register 注册消息
// 重复会 panic
func (p *Mgr) Register(messageID uint32, opts ...*options) {
	if p.messageMap == nil {
		p.messageMap = make(map[uint32]IMessage)
	}
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
	p.messageMap[messageID] = newDefaultMessage(opt)
}

func (p *Mgr) Find(messageID uint32) IMessage {
	if p.messageMap == nil {
		return nil
	}
	return p.messageMap[messageID]
}

// Replace 替换/覆盖(Override)
// 配置错误会 panic
func (p *Mgr) Replace(messageID uint32, opt *options) {
	if err := configure(opt); err != nil {
		xlog.PrintErr(xerror.Param, "%v messageID:%#x %v", xruntime.Location(), messageID, messageID)
		panic(errors.WithMessagef(xerror.Param, "%v messageID:%#x %v",
			xruntime.Location(), messageID, messageID))
	}
	p.messageMap[messageID] = newDefaultMessage(opt)
}
