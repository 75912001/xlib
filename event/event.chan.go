// Package event 提供了一个事件管理器，用于处理异步事件
// 使用 channel 作为事件通道
// 优点: 简单, 性能好
// 缺点: 需要预估事件通道容量, 如果容量不够, 会导致事件丢失
package event

import (
	"context"
	xconfigcommon "github.com/75912001/xlib/config/constants"
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xeventcommon "github.com/75912001/xlib/event/constants"
	xlog "github.com/75912001/xlib/log"
	xruntime "github.com/75912001/xlib/runtime"
	xutil "github.com/75912001/xlib/util"
	"runtime/debug"
	"time"
)

// ChanManager 事件管理器
type ChanManager struct {
	eventChan   chan any            // 事件通道
	onFunction  xcontrol.OnFunction // 事件处理器
	workerCount uint32              // 工作协程数量
	ctx         context.Context     // 上下文，用于控制协程生命周期
	cancel      context.CancelFunc
}

// NewChanManager 创建一个新的事件管理器
//
//	workerCount: 工作协程数量
//	handler: 事件处理函数
//	eventChanCapacity: 事件 chan 容量
func NewChanManager(workerCount uint32, onFunction xcontrol.OnFunction, eventChanCapacity uint32) *ChanManager {
	if workerCount <= 0 {
		workerCount = 1
	}
	if eventChanCapacity <= 0 {
		eventChanCapacity = xeventcommon.GEventChanCapacity
	}
	ctx, cancel := context.WithCancel(context.Background())

	return &ChanManager{
		eventChan:   make(chan any, eventChanCapacity),
		workerCount: workerCount,
		onFunction:  onFunction,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start 启动事件管理器
func (p *ChanManager) Start() {
	for range int(p.workerCount) {
		go p.worker()
	}
}

// Stop 停止事件管理器
func (p *ChanManager) Stop() {
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
func (p *ChanManager) Send(event any) {
	err := xutil.PushEventWithTimeout(p.eventChan, event, xconfigcommon.AddEventTimeoutDurationDefault)
	if err != nil {
		xlog.PrintErr(err)
	}
}

// worker 工作协程
func (p *ChanManager) worker() {
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
		case event := <-p.eventChan:
			if err := p.onFunction(event); err != nil {
				// 处理错误，可以选择记录日志或采取其他措施
				xlog.PrintErr(err)
			}
		}
	}
}

// 处理退出检查
func (p *ChanManager) handleQuitCheck() bool {
	cnt := len(p.eventChan)
	if cnt == 0 {
		xlog.PrintInfo("consume eventChan with length 0")
		return true
	}
	xlog.PrintfInfo("waiting for consume eventChan with length:%d", cnt)
	return false
}

func (p *ChanManager) Event(sb xcontrol.ISwitchButton, cb xcontrol.ICallBack) {
	p.Send(
		&xcontrol.Event{
			ISwitch:   sb,
			ICallBack: cb,
		},
	)
}
