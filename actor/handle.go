package actor

import (
	"fmt"

	xruntime "github.com/75912001/xlib/runtime"
)

func (p *Actor[KEY]) handleStop(event *BehaviorEvent) (response any, err error) {
	p.childrenMgr.Foreach(func(key KEY, child *Actor[KEY]) bool { // 停止所有子 Actor
		if event.IsSync() { // 同步
			_, _ = child.SendEventWithResponse(event)
		} else { // 异步
			child.SendEvent(event)
		}
		return true
	})
	p.childrenMgr.Clear()
	if event.IsSync() { // 响应调用方（如果是同步调用）
		event.responseChan <- &behaviorResponse{data: nil, err: nil}
	}
	// 停止自己的事件管理器
	p.eventMgr.Stop()
	return nil, nil
}

func (p *Actor[KEY]) handleRemoveChild(event *BehaviorEvent) (response any, err error) {
	childKey := event.Args[0]

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
	key := event.Args[0].(KEY)
	behavior := event.Args[1].(Behavior)

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
	key := event.Args[0].(KEY)

	// 检查是否已存在
	child := p.childrenMgr.Get(key)
	if child == nil {
		return nil, fmt.Errorf("child %v is not exist %v", key, xruntime.Location())
	}
	return any(child), nil
}
