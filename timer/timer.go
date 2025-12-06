// Package 定时器
// 优先级: 到期时间,加入顺序

package timer

import (
	"container/heap"
	"container/list"
	"context"
	xconfig "github.com/75912001/xlib/config"
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xruntime "github.com/75912001/xlib/runtime"
	xtimerconstants "github.com/75912001/xlib/timer/constants"
	"runtime/debug"
	"sync"
	"time"
)

// 定时器
type defaultTimer struct {
	secondSlice     [cycleSize]list.List // 时间轮-数组. 秒,数据
	millisecondList list.List            // 毫秒级-list
	milliTaskHeap   *MillisecondMinHeap  // 毫秒级-小顶堆

	cancelFunc      context.CancelFunc
	waitGroup       sync.WaitGroup // Stop 等待信号
	milliSecondChan chan any       // 毫秒, channel
	secondChan      chan any       // 秒, channel

	// 统计-秒-定时器-数量
	secondCount uint64
	// 统计-毫秒-定时器-数量
	millisecondCount uint64
}

func NewTimer() ITimer {
	return &defaultTimer{
		milliTaskHeap: InitMilliTaskHeap(),
	}
}

// 每秒更新
func (p *defaultTimer) funcSecond(ctx context.Context) {
	defer func() {
		if xruntime.IsRelease() {
			if err := recover(); err != nil {
				xlog.PrintErr(xerror.GoroutinePanic, err, string(debug.Stack()))
			}
		}
		p.waitGroup.Done()
		xlog.PrintInfo(xerror.GoroutineDone)
	}()
	scanSecondDuration := xconfig.GConfigMgr.Timer.GetScanSecondDuration()
	idleDelay := time.NewTimer(scanSecondDuration)
	defer func() {
		idleDelay.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			xlog.PrintInfo(xerror.GoroutineDone)
			return
		case v := <-p.secondChan:
			s := v.(*Second)
			duration := s.expire - ShadowTimestamp()
			if duration < 0 { // 到期
				duration = 0
			}
			cycleIdx := searchCycleIdx(duration)
			p.pushBackCycle(s, cycleIdx)
			p.secondCount++
			sortByExpire(&p.secondSlice[cycleIdx])
		case <-idleDelay.C:
			idleDelay.Reset(scanSecondDuration)
			p.scanSecond(ShadowTimestamp())
		}
	}
}

// 每 Millisecond 个毫秒更新
func (p *defaultTimer) funcMillisecond(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			xlog.PrintErr(xerror.GoroutinePanic, err, string(debug.Stack()))
		}
		p.waitGroup.Done()
		xlog.PrintInfo(xerror.GoroutineDone)
	}()
	scanMillisecondDuration := xconfig.GConfigMgr.Timer.GetScanMillisecondDuration()
	scanMillisecond := scanMillisecondDuration / time.Millisecond
	idleDelay := time.NewTimer(scanMillisecondDuration)
	defer func() {
		idleDelay.Stop()
	}()
	nextMillisecond := time.Duration(time.Now().UnixMilli()) + scanMillisecond

	for {
		select {
		case <-ctx.Done():
			xlog.PrintInfo(xerror.GoroutineDone)
			return
		case v := <-p.milliSecondChan:
			millisecond := v.(*Millisecond)
			p.millisecondCount++
			switch xconfig.GConfigMgr.Timer.GetMillisecondType() {
			case xtimerconstants.MillisecondTypeList:
				p.millisecondList.PushBack(millisecond)
				sortByExpire(&p.millisecondList)
			case xtimerconstants.MillisecondTypeMinHeap:
				heap.Push(p.milliTaskHeap, NewMilliTask(millisecond.expire, millisecond))
			}
		case <-idleDelay.C:
			nowMillisecond := time.Now().UnixMilli()
			reset := scanMillisecondDuration - (time.Duration(nowMillisecond)-nextMillisecond)*time.Millisecond
			idleDelay.Reset(reset)

			nextMillisecond += scanMillisecond
			p.scanMillisecond(nowMillisecond)
		}
	}
}

// Start
func (p *defaultTimer) Start(ctx context.Context) error {
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	p.cancelFunc = cancelFunc

	{
		p.secondChan = make(chan any, 1000)
		p.waitGroup.Add(1)

		go p.funcSecond(ctxWithCancel)
	}
	{
		p.milliSecondChan = make(chan any, 1000)
		p.waitGroup.Add(1)

		go p.funcMillisecond(ctxWithCancel)
	}
	return nil
}

// Stop 停止服务
func (p *defaultTimer) Stop() {
	if p.cancelFunc != nil {
		p.cancelFunc()
		// 等待 Second, milliSecond goroutine退出.
		p.waitGroup.Wait()
		p.cancelFunc = nil
	}
}

// AddMillisecond 添加毫秒级定时器
//
//	参数:
//		callBackFunc: 回调接口
//		expireMillisecond: 过期毫秒数
//	返回值:
//		毫秒定时器
func (p *defaultTimer) AddMillisecond(callBackFunc xcontrol.ICallBack, expireMillisecond int64, out xcontrol.IOut) *Millisecond {
	t := &Millisecond{
		ICallBack:     callBackFunc,
		ISwitchButton: xcontrol.NewSwitchButton(true),
		expire:        expireMillisecond,
		IOut:          out,
	}
	p.milliSecondChan <- t
	return t
}

// DelMillisecond 删除毫秒级定时器
//
//	参数:
//		毫秒定时器
func (p *defaultTimer) DelMillisecond(t *Millisecond) {
	t.Delete()
}

// 扫描毫秒级定时器
//
//	参数:
//		ms: 到期毫秒数
func (p *defaultTimer) scanMillisecond(ms int64) {
	switch xconfig.GConfigMgr.Timer.GetMillisecondType() {
	case xtimerconstants.MillisecondTypeList:
		var next *list.Element
		for e := p.millisecondList.Front(); e != nil; e = next {
			t := e.Value.(*Millisecond)
			if t.ISwitchButton.IsOff() {
				next = e.Next()
				p.millisecondList.Remove(e)
				p.millisecondCount--
				continue
			}
			if t.expire <= ms {
				t.IOut.Send(
					&xcontrol.Event{
						ISwitch:   t.ISwitchButton,
						ICallBack: t.ICallBack,
					},
				)
				next = e.Next()
				p.millisecondList.Remove(e)
				p.millisecondCount--
				continue
			}
			break
		}
	case xtimerconstants.MillisecondTypeMinHeap:
		for p.milliTaskHeap.Len() > 0 {
			milliTask := (*p.milliTaskHeap)[0]               // 只看堆顶
			if milliTask.millisecond.ISwitchButton.IsOff() { // 已删除
				heap.Pop(p.milliTaskHeap) // 弹出任务
				p.millisecondCount--
				continue
			}
			if ms < milliTask.expire {
				break // 堆顶未到期，后面都不会到期
			}
			heap.Pop(p.milliTaskHeap) // 弹出到期任务
			milliTask.millisecond.IOut.Send(
				&xcontrol.Event{
					ISwitch:   milliTask.millisecond.ISwitchButton,
					ICallBack: milliTask.millisecond.ICallBack,
				},
			)
			p.millisecondCount--
		}
	}

}

// AddSecond 添加秒级定时器
//
//	参数:
//		callBackFunc: 回调接口
//		expire: 过期秒数
//	返回值:
//		秒定时器
func (p *defaultTimer) AddSecond(callBackFunc xcontrol.ICallBack, expire int64, out xcontrol.IOut) *Second {
	t := &Second{
		Millisecond: &Millisecond{
			ISwitchButton: xcontrol.NewSwitchButton(true),
			ICallBack:     callBackFunc,
			expire:        expire,
			IOut:          out,
		},
	}
	p.secondChan <- t
	return t
}

// DelSecond 删除秒级定时器
func (p *defaultTimer) DelSecond(t *Second) {
	t.Delete()
}

// 将秒级定时器,添加到轮转IDX的末尾.
//
//	参数:
//		timerSecond: 秒定时器
//		cycleIdx: 轮序号
func (p *defaultTimer) pushBackCycle(timerSecond *Second, cycleIdx int) {
	l := &p.secondSlice[cycleIdx]
	l.PushBack(timerSecond)
	if 100000 < l.Len() { // 监控告警
		xlog.GLog.Warnf("time wheel slot %v overload, len %v", cycleIdx, l.Len())
	}
}

// 扫描秒级定时器
//
//	timestamp: 到期时间戳
func (p *defaultTimer) scanSecond(timestamp int64) {
	var next *list.Element
	cycle0 := &p.secondSlice[0]
	for e := cycle0.Front(); nil != e; e = next {
		t := e.Value.(*Second)
		if t.ISwitchButton.IsOff() {
			next = e.Next()
			cycle0.Remove(e)
			p.secondCount--
			continue
		}
		if t.expire <= timestamp {
			t.IOut.Send(
				&xcontrol.Event{
					ISwitch:   t.ISwitchButton,
					ICallBack: t.ICallBack,
				},
			)
			next = e.Next()
			cycle0.Remove(e)
			p.secondCount--
			continue
		}
		break
	}
	if cycle0.Len() != 0 { // 如果当前的 cycle 中还有元素,则不需要之后的cycle向前移动
		return
	}
	// 更新时间轮,从序号为1的数组开始
	for idx := 1; idx < cycleSize; idx++ {
		c := &p.secondSlice[idx]
		for e := c.Front(); e != nil; e = next {
			t := e.Value.(*Second)
			if t.ISwitchButton.IsOff() {
				next = e.Next()
				c.Remove(e)
				p.secondCount--
				continue
			}
			if t.expire <= timestamp {
				t.IOut.Send(
					&xcontrol.Event{
						ISwitch:   t.ISwitchButton,
						ICallBack: t.ICallBack,
					},
				)
				next = e.Next()
				c.Remove(e)
				p.secondCount--
				continue
			}
			if newIdx := findPrevCycleIdx(t.expire-timestamp, idx); idx != newIdx {
				next = e.Next()
				c.Remove(e)
				p.pushBackCycle(t, newIdx)
				continue
			}
			break
		}
		if c.Len() != 0 { // 如果当前的 cycle 中还有元素,则不需要之后的cycle向前移动
			break
		}
	}
}
