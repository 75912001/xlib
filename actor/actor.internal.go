package actor

import (
	"fmt"

	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

func (p *Actor[KEY]) stop(msg *Msg) (resp any, err error) {
	p.childMgr.Foreach(func(key KEY, child *Actor[KEY]) bool { // 停止所有子 Actor
		if msg.IsSync() { // 同步
			_, _ = child.SendMsgAsync(
				NewMsg(msg.Ctx, SystemReservedCommand_Stop),
			)
		} else { // 异步
			child.SendMsg(
				NewMsg(msg.Ctx, SystemReservedCommand_Stop),
			)
		}
		return true
	})
	p.childMgr.Clear()
	// 停止自己的消息管理器
	p.msgMgr.Stop()
	return nil, nil
}

func (p *Actor[KEY]) removeChild(msg *Msg) (resp any, err error) {
	if len(msg.Args) < 1 {
		return nil, errors.WithMessage(xerror.ParamCountNotMatch, xruntime.Location())
	}
	childKey, ok := msg.Args[0].(KEY)
	if !ok {
		return nil, fmt.Errorf("invalid child key type %v", xruntime.Location())
	}
	child, ok := p.childMgr.Find(childKey)
	if !ok {
		return nil, fmt.Errorf("child not exist %v %v", childKey, xruntime.Location())
	}
	if msg.IsSync() {
		_, _ = child.SendMsgAsync(
			NewMsg(msg.Ctx, SystemReservedCommand_Stop),
		)
	} else {
		child.SendMsg(
			NewMsg(msg.Ctx, SystemReservedCommand_Stop),
		)
	}
	p.childMgr.Del(childKey)
	return nil, nil
}

func (p *Actor[KEY]) spawn(msg *Msg) (resp any, err error) {
	if len(msg.Args) < 2 {
		return nil, errors.WithMessage(xerror.ParamCountNotMatch, xruntime.Location())
	}
	key, ok := msg.Args[0].(KEY)
	if !ok {
		return nil, fmt.Errorf("invalid child key type %v", xruntime.Location())
	}
	behavior, ok := msg.Args[1].(Behavior)
	if !ok {
		return nil, fmt.Errorf("invalid child behavior type %v", xruntime.Location())
	}
	// 检查是否已存在
	if p.childMgr.Get(key) != nil {
		return nil, fmt.Errorf("child %v is already exist %v", key, xruntime.Location())
	}
	// 创建并启动子 Actor
	child := NewActor(key, p, behavior)
	p.childMgr.Add(key, child)
	child.Start()
	return any(child), nil
}

func (p *Actor[KEY]) getChild(msg *Msg) (resp any, err error) {
	if len(msg.Args) < 1 {
		return nil, errors.WithMessage(xerror.ParamCountNotMatch, xruntime.Location())
	}
	key, ok := msg.Args[0].(KEY)
	if !ok {
		return nil, fmt.Errorf("invalid child key type %v", xruntime.Location())
	}
	// 检查是否已存在
	child := p.childMgr.Get(key)
	if child == nil {
		return nil, fmt.Errorf("child %v is not exist %v", key, xruntime.Location())
	}
	return any(child), nil
}
