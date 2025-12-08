// Package event 提供了一个事件管理器，用于处理异步事件
// 使用 list.List 作为事件通道
// 优点: 不需要预估事件通道容量, 不会丢失事件
// 缺点: 性能不如 channel
package event

import (
	"container/list"
	"context"
	"runtime/debug"
	"sync"
	"time"

	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xruntime "github.com/75912001/xlib/runtime"
)

type Manager ListMgr

// ListMgr 事件管理器
type ListMgr struct {
	queueMu    sync.Mutex // 保护 events
	events     *list.List
	notifyChan chan struct{} // 通知有新消息

	onFunction  xcontrol.OnFunction // 事件处理器
	workerCount uint32              // 工作协程数量
	ctx         context.Context     // 上下文，用于控制协程生命周期
	cancel      context.CancelFunc
}

// NewListMgr 创建一个新的事件管理器
//
//	workerCount: 工作协程数量
//	handler: 事件处理函数
func NewListMgr(workerCount uint32, onFunction xcontrol.OnFunction) *ListMgr {
	if workerCount <= 0 {
		workerCount = 1
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &ListMgr{
		events:      list.New(),
		notifyChan:  make(chan struct{}, int(workerCount)),
		workerCount: workerCount,
		onFunction:  onFunction,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start 启动事件管理器
func (p *ListMgr) Start() {
	for range int(p.workerCount) {
		go p.worker()
	}
}

// Stop 停止事件管理器
func (p *ListMgr) Stop() {
	go func() { // 检查协程是否可以结束
		idleDuration := 100 * time.Millisecond
		idleDelay := time.NewTimer(idleDuration)
		defer func() {
			idleDelay.Stop()
		}()
		for range idleDelay.C {
			idleDelay.Reset(idleDuration)
			if p.handleQuitCheck() {
				p.cancel()
				return
			}
		}
	}()
}

// Send 发送事件到管理器
func (p *ListMgr) Send(events ...any) {
	p.queueMu.Lock()
	for _, event := range events {
		p.events.PushBack(event)
	}
	p.queueMu.Unlock()
	select {
	case p.notifyChan <- struct{}{}:
	default: // channel 已满，说明所有 worker 都在工作
	}
}

// worker 工作协程
func (p *ListMgr) worker() {
	defer func() {
		if xruntime.IsRelease() {
			if r := recover(); r != nil {
				xlog.PrintfErr("worker panic: %v", r)
				// 打印堆栈
				debug.PrintStack()
			}
		}
		xlog.PrintInfo(xerror.GoroutineDone.Error())
	}()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-p.notifyChan:
			for {
				p.queueMu.Lock()
				element := p.events.Front()
				if element == nil { // 没有事件
					p.queueMu.Unlock()
					break
				}
				event := element.Value
				p.events.Remove(element)
				p.queueMu.Unlock()

				if err := p.onFunction(event); err != nil {
					// 处理错误，可以选择记录日志或采取其他措施
					xlog.GLog.Errorf("event:%+v err:%v", event, err)
				}
			}
		}
	}
}

// 处理退出检查
func (p *ListMgr) handleQuitCheck() bool {
	p.queueMu.Lock()
	cnt := p.events.Len()
	p.queueMu.Unlock()
	if cnt == 0 {
		xlog.PrintInfo("consume eventChan with length 0")
		return true
	}
	xlog.PrintfInfo("waiting for consume eventChan with length:%d", cnt)
	return false
}
