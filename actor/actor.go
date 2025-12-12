package actor

import (
	"time"

	xerror "github.com/75912001/xlib/error"
	xevent "github.com/75912001/xlib/event"
	xmap "github.com/75912001/xlib/map"
	xpool "github.com/75912001/xlib/pool"
	xruntime "github.com/75912001/xlib/runtime"
	xstatistics "github.com/75912001/xlib/statistics"
	"github.com/pkg/errors"
)

// Actor 表示一个 actor 实例
type Actor[KEY comparable] struct {
	key      KEY             // actor 的唯一标识
	msgMgr   *xevent.ListMgr // 消息-管理器
	behavior Behavior        // 行为函数
	parent   *Actor[KEY]     // 父 actor, 可以为 nil

	childMgr *xmap.MapMgr[KEY, *Actor[KEY]] // 子 actor 管理器

	Statistics *xstatistics.Statistics // 统计数据
}

// NewActor 创建一个新的 actor
func NewActor[KEY comparable](key KEY, parent *Actor[KEY], behavior Behavior) *Actor[KEY] {
	actor := &Actor[KEY]{
		key:        key,
		behavior:   behavior,
		parent:     parent,
		childMgr:   xmap.NewMapMgr[KEY, *Actor[KEY]](),
		Statistics: xstatistics.NewStatistics(),
	}
	actor.msgMgr = xevent.NewListMgr(1, actor.process)
	return actor
}

func (p *Actor[KEY]) GetKey() KEY {
	return p.key
}

// GetParent 获取父 actor
func (p *Actor[KEY]) GetParent() *Actor[KEY] {
	return p.parent
}

// Start 启动actor
func (p *Actor[KEY]) Start() {
	p.msgMgr.Start()
}

// Send 发送消息到 actor (异步)
func (p *Actor[KEY]) Send(messages ...any) {
	p.msgMgr.Send(messages...)
}

// SendMsg 发送消息到 actor (异步)
func (p *Actor[KEY]) SendMsg(messages ...*Msg) {
	// 将每个元素转换为any类型
	messagesAny := make([]any, len(messages))
	for i, msg := range messages {
		messagesAny[i] = msg
	}
	p.msgMgr.Send(messagesAny...)
}

// SendMsgAsync 发送消息到 actor, 并等待响应 (同步)
//
//	参数:
//		msg: 行为-消息
//			Ctx: 上下文, 可用于设置超时, 如果 ctx 没有设置超时, 则使用默认 60 秒超时
//	返回:
//		resp: 响应数据
//		err: 错误
func (p *Actor[KEY]) SendMsgAsync(msg *Msg) (resp any, err error) {
	msg.withAsyncChan(make(chan *behaviorResponse, 1))
	p.SendMsg(msg)

	// 判断 context 是否设置了超时
	var timer *time.Timer
	var hasDeadline bool
	if msg.Ctx != nil {
		_, hasDeadline = msg.Ctx.Deadline()
	}
	if msg.Ctx != nil && hasDeadline { // 设置了 context, 设置超时 => 使用 context 的超时
		select {
		case res := <-msg.asyncChan:
			return res.respData, res.err
		case <-msg.Ctx.Done():
			return nil, errors.WithMessagef(xerror.Timeout, "SendMsgAsync context timeout. msg:%v %v", msg, xruntime.Location())
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
	case res := <-msg.asyncChan:
		return res.respData, res.err
	case <-timer.C:
		return nil, errors.WithMessagef(xerror.Timeout, "SendMsgAsync default timeout after 60 seconds. event:%v %v", msg, xruntime.Location())
	}
}
