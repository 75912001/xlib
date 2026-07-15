package actor

import (
	"sync"
	"sync/atomic"
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
	key         KEY // actor 的唯一标识
	state       atomic.Uint32
	lifecycleMu sync.RWMutex
	msgMgr      *xevent.ListMgr // 消息-管理器
	behavior    Behavior        // 行为函数
	parent      *Actor[KEY]     // 父 actor, 可以为 nil

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
	actor.state.Store(StateNew)
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
func (p *Actor[KEY]) Start() error {
	if !p.state.CompareAndSwap(
		StateNew,
		StateRunning,
	) {
		return errors.WithMessagef(xerror.Fail, "%v", xruntime.Location())
	}

	p.msgMgr.Start()
	return nil
}

// Send 发送消息到 actor (异步)
func (p *Actor[KEY]) Send(messages ...any) error {
	p.lifecycleMu.RLock()
	defer p.lifecycleMu.RUnlock()

	if p.state.Load() != StateRunning {
		return errors.WithMessagef(xerror.Fail, "%v", xruntime.Location())
	}

	p.msgMgr.Send(messages...)
	return nil
}

// SendMsg 发送消息到 actor (异步)
func (p *Actor[KEY]) SendMsg(messages ...*Msg) error {
	p.lifecycleMu.RLock()
	defer p.lifecycleMu.RUnlock()

	if p.state.Load() != StateRunning {
		return errors.WithMessagef(xerror.Fail, "%v", xruntime.Location())
	}
	// 将每个元素转换为any类型
	messagesAny := make([]any, len(messages))
	for i, msg := range messages {
		messagesAny[i] = msg
	}
	p.msgMgr.Send(messagesAny...)
	return nil
}

// SendMsgSync 发送消息到 actor, 并等待响应 (同步)
//
//	参数:
//		msg: 行为-消息
//			Ctx: 上下文, 可用于设置超时, 如果 ctx 没有设置超时, 则使用默认 60 秒超时
//	返回:
//		resp: 响应数据
//		err: 错误
func (p *Actor[KEY]) SendMsgSync(msg *Msg) (resp any, err error) {
	msg.withSyncChan(make(chan *behaviorResponse, 1))
	err = p.SendMsg(msg)
	if err != nil {
		return nil, errors.WithMessagef(err, "SendMsgSync send msg error. msg:%v %v", msg, xruntime.Location())
	}

	var hasDeadline bool
	var ctxDone <-chan struct{}
	if msg.Ctx != nil {
		_, hasDeadline = msg.Ctx.Deadline()
		ctxDone = msg.Ctx.Done()
	}

	if hasDeadline { // 设置了超时, 使用 context 的超时, 不加 60 秒上限
		select {
		case res := <-msg.syncChan:
			return res.respData, res.err
		case <-ctxDone:
			return nil, errors.WithMessagef(xerror.Timeout, "SendMsgSync context timeout. msg:%v %v", msg, xruntime.Location())
		}
	}
	// 未设置超时, 使用 60 秒默认超时兜底; ctxDone 为 nil 时永不触发 (nil channel)
	const timeout = 60 * time.Second
	timer := xpool.Timer.Get()
	timer.Reset(timeout)
	defer func() {
		xpool.Timer.Put(timer)
	}()
	select {
	case res := <-msg.syncChan:
		return res.respData, res.err
	case <-ctxDone:
		return nil, errors.WithMessagef(xerror.Timeout, "SendMsgSync context timeout. msg:%v %v", msg, xruntime.Location())
	case <-timer.C:
		return nil, errors.WithMessagef(xerror.Timeout, "SendMsgSync default timeout after 60 seconds. event:%v %v", msg, xruntime.Location())
	}
}
