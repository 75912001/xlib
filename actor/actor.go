package actor

import (
	xerror "github.com/75912001/xlib/error"
	xevent "github.com/75912001/xlib/event"
	xmap "github.com/75912001/xlib/map"
	xpool "github.com/75912001/xlib/pool"
	xruntime "github.com/75912001/xlib/runtime"
	xstatistics "github.com/75912001/xlib/statistics"
	"github.com/pkg/errors"
	"sync"
	"time"
)

// Actor 表示一个 actor 实例
type Actor[KEY comparable] struct {
	key      KEY                 // actor 的唯一标识
	eventMgr *xevent.ListManager // 事件管理器
	behavior Behavior            // 行为函数
	parent   *Actor[KEY]         // 父 actor, 可以为 nil

	childrenMgr *xmap.MapMgr[KEY, *Actor[KEY]] // 子 actor 管理器
	childrenMu  sync.RWMutex                   // childrenMgr mutex

	Statistics *xstatistics.Statistics // 统计数据
}

// NewActor 创建一个新的 actor
func NewActor[KEY comparable](key KEY, parent *Actor[KEY], behavior Behavior) *Actor[KEY] {
	actor := &Actor[KEY]{
		key:         key,
		behavior:    behavior,
		parent:      parent,
		childrenMgr: xmap.NewMapMgr[KEY, *Actor[KEY]](),
		Statistics:  xstatistics.NewStatistics(),
	}
	actor.eventMgr = xevent.NewListManager(1, actor.process)
	return actor
}

func (p *Actor[KEY]) GetKey() KEY {
	return p.key
}

// Start 启动actor
func (p *Actor[KEY]) Start() {
	p.eventMgr.Start()
}

// Stop 停止 actor
func (p *Actor[KEY]) Stop() {
	p.childrenMu.Lock()
	// 停止所有子 actor
	p.childrenMgr.Foreach(func(key KEY, child *Actor[KEY]) bool {
		child.Stop()
		return true
	})
	p.childrenMgr.Clear()
	p.childrenMu.Unlock()

	// 停止自己
	p.eventMgr.Stop()
}

// Send 发送消息到 actor (异步)
func (p *Actor[KEY]) Send(events ...any) {
	p.eventMgr.Send(events...)
}

// SendEvent 发送到 actor (异步)
func (p *Actor[KEY]) SendEvent(events ...*BehaviorEvent) {
	// 将每个元素转换为any类型
	eventsAny := make([]any, len(events))
	for i, event := range events {
		eventsAny[i] = event
	}
	p.eventMgr.Send(eventsAny...)
}

// SendEventWithResponse 发送消息到 actor, 并等待响应 (同步)
//
//	参数:
//		event: 行为事件
//			Ctx: 上下文, 可用于设置超时, 如果 ctx 没有设置超时, 则使用默认 60 秒超时
//	返回:
//		response: 响应数据
//		err: 错误
func (p *Actor[KEY]) SendEventWithResponse(event *BehaviorEvent) (response any, err error) {
	event.withResponseChan(make(chan *behaviorResponse, 1))
	p.SendEvent(event)

	// 判断 context 是否设置了超时
	var timer *time.Timer
	var hasDeadline bool
	if event.Ctx != nil {
		_, hasDeadline = event.Ctx.Deadline()
	}
	if event.Ctx != nil && hasDeadline { // 设置了 context, 设置超时 => 使用 context 的超时
		select {
		case res := <-event.responseChan:
			return res.data, res.err
		case <-event.Ctx.Done():
			return nil, errors.WithMessagef(xerror.Timeout, "SendEventWithResponse context timeout. event:%v %v", event, xruntime.Location())
		}
	}
	// 没有设置 context, 没有设置超时 => 使用默认 60 秒超时
	const timeout = 60 * time.Second
	timer = xpool.Timer.Get()
	ok := timer.Reset(timeout)
	if !ok {
		xpool.Timer.Put(timer)
		timer = time.NewTimer(timeout)
		defer timer.Stop()
	} else {
		defer func() {
			timer.Stop()
			xpool.Timer.Put(timer)
		}()
	}
	select {
	case res := <-event.responseChan:
		return res.data, res.err
	case <-timer.C:
		return nil, errors.WithMessagef(xerror.Timeout, "SendEventWithResponse default timeout after 60 seconds. event:%v %v", event, xruntime.Location())
	}
}

// process 处理消息的主循环
func (p *Actor[KEY]) process(args ...any) (err error) {
	start := time.Now()
	defer func() {
		p.Statistics.ProcessTime += time.Since(start)
		p.Statistics.Count++
		if err != nil {
			p.Statistics.ErrorCount++
		}
	}()
	msg := args[0]
	var response any
	switch event := msg.(type) {
	case *BehaviorEvent: // 处理 BehaviorEvent
		if p.behavior == nil {
			err = errors.WithMessagef(xerror.NoBehavior, "actor has no behavior. key:%v %v",
				p.key, xruntime.Location())
			if event.IsSync() { // 同步调用
				event.responseChan <- &behaviorResponse{
					data: response,
					err:  err,
				}
			}
			return
		}
		p.behavior, response, err = p.behavior(event)
		if err != nil {
			err = errors.WithMessagef(err, "actor process error. %v", xruntime.Location())
		}
		if event.IsSync() { // 同步调用
			event.responseChan <- &behaviorResponse{
				data: response,
				err:  err,
			}
		}
		return
	default: // 类型未知
		if p.behavior == nil {
			err = errors.WithMessagef(xerror.NoBehavior, "actor has no behavior. %v", xruntime.Location())
			return
		}
		p.behavior, response, err = p.behavior(msg)
		_ = response
		err = errors.WithMessagef(err, "actor process error. %v", xruntime.Location())
		return
	}
}

// Spawn 创建子 actor
func (p *Actor[KEY]) Spawn(key KEY, behavior Behavior) *Actor[KEY] {
	p.childrenMu.Lock()
	defer p.childrenMu.Unlock()

	child := NewActor(key, p, behavior)
	p.childrenMgr.Add(key, child)
	child.Start()
	return child
}

// GetChild 获取子 actor
func (p *Actor[KEY]) GetChild(key KEY) *Actor[KEY] {
	p.childrenMu.RLock()
	defer p.childrenMu.RUnlock()

	return p.childrenMgr.Get(key)
}

// RemoveChild 移除子 actor
func (p *Actor[KEY]) RemoveChild(key KEY) {
	p.childrenMu.Lock()
	defer p.childrenMu.Unlock()

	child, ok := p.childrenMgr.Find(key)
	if !ok {
		return
	}
	child.Stop()
	p.childrenMgr.Del(key)
}

// GetParent 获取父 actor
func (p *Actor[KEY]) GetParent() *Actor[KEY] {
	return p.parent
}
