package actor

import (
	"fmt"

	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

func (p *Actor[KEY]) handleStop(event *BehaviorEvent) (response any, err error) {
	p.childrenMgr.Foreach(func(key KEY, child *Actor[KEY]) bool { // 停止所有子 Actor
		if event.IsSync() { // 同步
			_, _ = child.SendEventWithResponse(
				NewBehaviorEvent(event.Ctx, SystemReservedCommand_Stop),
			)
		} else { // 异步
			child.SendEvent(
				NewBehaviorEvent(event.Ctx, SystemReservedCommand_Stop),
			)
		}
		return true
	})
	p.childrenMgr.Clear()
	// 停止自己的消息管理器
	p.messageMgr.Stop()
	return nil, nil
}

func (p *Actor[KEY]) handleRemoveChild(event *BehaviorEvent) (response any, err error) {
	if len(event.Args) < 1 {
		return nil, errors.WithMessage(xerror.ParamCountNotMatch, xruntime.Location())
	}
	childKey, ok := event.Args[0].(KEY)
	if !ok {
		return nil, fmt.Errorf("invalid child key type %v", xruntime.Location())
	}
	child, ok := p.childrenMgr.Find(childKey)
	if !ok {
		return nil, fmt.Errorf("child not exist %v %v", childKey, xruntime.Location())
	}
	if event.IsSync() {
		_, _ = child.SendEventWithResponse(
			NewBehaviorEvent(event.Ctx, SystemReservedCommand_Stop),
		)
	} else {
		child.SendEvent(
			NewBehaviorEvent(event.Ctx, SystemReservedCommand_Stop),
		)
	}
	p.childrenMgr.Del(childKey)
	return nil, nil
}

func (p *Actor[KEY]) handleSpawn(event *BehaviorEvent) (response any, err error) {
	if len(event.Args) < 2 {
		return nil, errors.WithMessage(xerror.ParamCountNotMatch, xruntime.Location())
	}
	key, ok := event.Args[0].(KEY)
	if !ok {
		return nil, fmt.Errorf("invalid child key type %v", xruntime.Location())
	}
	behavior, ok := event.Args[1].(Behavior)
	if !ok {
		return nil, fmt.Errorf("invalid child behavior type %v", xruntime.Location())
	}
	// 检查是否已存在
	if p.childrenMgr.Get(key) != nil {
		return nil, fmt.Errorf("child %v is already exist %v", key, xruntime.Location())
	}
	// 创建并启动子 Actor
	child := NewActor(key, p, behavior)
	p.childrenMgr.Add(key, child)
	child.Start()
	return any(child), nil
}

func (p *Actor[KEY]) handleGetChild(event *BehaviorEvent) (response any, err error) {
	if len(event.Args) < 1 {
		return nil, errors.WithMessage(xerror.ParamCountNotMatch, xruntime.Location())
	}
	key, ok := event.Args[0].(KEY)
	if !ok {
		return nil, fmt.Errorf("invalid child key type %v", xruntime.Location())
	}
	// 检查是否已存在
	child := p.childrenMgr.Get(key)
	if child == nil {
		return nil, fmt.Errorf("child %v is not exist %v", key, xruntime.Location())
	}
	return any(child), nil
}
