// Package 定时器
// 优先级: 到期时间,加入顺序

package timer

import (
	"container/list"
	"context"
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"runtime/debug"
	"sync"
	"time"
)

// 定时器
type defaultTimer struct {
	opts            *options
	secondSlice     [cycleSize]list.List // 时间轮-数组. 秒,数据
	millisecondList list.List            // list. 毫秒,数据
	cancelFunc      context.CancelFunc
	waitGroup       sync.WaitGroup   // Stop 等待信号
	milliSecondChan chan interface{} // 毫秒, channel
	secondChan      chan interface{} // 秒, channel
}

func NewTimer() ITimer {
	return &defaultTimer{}
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
	idleDelay := time.NewTimer(*p.opts.scanSecondDuration)
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
			if duration < 0 {
				duration = 0
			}
			p.pushBackCycle(s, searchCycleIdxIteration(duration), true)
		case <-idleDelay.C:
			idleDelay.Reset(*p.opts.scanSecondDuration)
			p.scanSecond(ShadowTimestamp())
		}
	}
}

// 每 millisecond 个毫秒更新
func (p *defaultTimer) funcMillisecond(ctx context.Context) {
	defer func() {
		if xruntime.IsRelease() {
			if err := recover(); err != nil {
				xlog.PrintErr(xerror.GoroutinePanic, err, string(debug.Stack()))
			}
		}
		p.waitGroup.Done()
		xlog.PrintInfo(xerror.GoroutineDone)
	}()
	scanMillisecondDuration := *p.opts.scanMillisecondDuration
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
			p.millisecondList.PushBack(v)
			moveLastElementToProperPosition(&p.millisecondList)
		case <-idleDelay.C:
			nowMillisecond := time.Now().UnixMilli()
			reset := scanMillisecondDuration - (time.Duration(nowMillisecond)-nextMillisecond)*time.Millisecond
			idleDelay.Reset(reset)

			nextMillisecond += scanMillisecond
			p.scanMillisecond(nowMillisecond)
		}
	}
}

// 移动最后一个元素到合适的位置,移动到大于他的元素的前面[实现按照时间排序,加入顺序排序]
// e.g.: 1,2,2,3,4,4,3 => 1,2,2,3,3,4,4 [将最后一个元素移动到4的前面]
// todo menglc 可以优化为,二分查找,然后插入
func moveLastElementToProperPosition(l *list.List) {
	lastElement := l.Back() // 获取最后一个元素
	target := lastElement.Value.(*millisecond)
	var element *list.Element
	for element = lastElement.Prev(); element != nil; element = element.Prev() {
		current := element.Value.(*millisecond)
		if current.expire <= target.expire {
			l.MoveAfter(lastElement, element)
			return
		}
	}
	if element == nil {
		// 如果没有找到比目标小或等于的元素，将目标元素移动到列表的前面
		l.MoveToFront(lastElement)
	}
}

// Start
// [NOTE] 处理定时器相关数据,必须与该 outgoingTimeoutChan 线性处理.如:在同一个 goroutine select 中处理数据
func (p *defaultTimer) Start(ctx context.Context, opts ...*options) error {
	p.opts = &options{}
	p.opts = p.opts.merge(opts...)
	if err := p.opts.configure(); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}

	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	p.cancelFunc = cancelFunc

	if p.opts.scanSecondDuration != nil {
		p.secondChan = make(chan interface{}, 100)
		p.waitGroup.Add(1)

		go p.funcSecond(ctxWithCancel)
	}
	if p.opts.scanMillisecondDuration != nil {
		p.milliSecondChan = make(chan interface{}, 100)
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
func (p *defaultTimer) AddMillisecond(callBackFunc xcontrol.ICallBack, expireMillisecond int64) *millisecond {
	t := &millisecond{
		ICallBack: callBackFunc,
		ISwitch:   xcontrol.NewSwitchButton(true),
		expire:    expireMillisecond,
	}
	p.milliSecondChan <- t
	return t
}

// DelMillisecond 删除毫秒级定时器
//
//	[NOTE] 必须与该 outgoingTimeoutChan 线性处理.如:在同一个 goroutine select 中处理数据
//	参数:
//		毫秒定时器
func (p *defaultTimer) DelMillisecond(t *millisecond) {
	t.reset()
}

// 扫描毫秒级定时器
//
//	参数:
//		ms: 到期毫秒数
func (p *defaultTimer) scanMillisecond(ms int64) {
	var next *list.Element
	for e := p.millisecondList.Front(); e != nil; e = next {
		t := e.Value.(*millisecond)
		if t.ISwitch.IsOff() {
			next = e.Next()
			p.millisecondList.Remove(e)
			continue
		}
		if t.expire <= ms {
			p.opts.outgoingTimeoutChan <- &EventTimerMillisecond{
				ISwitch:   t.ISwitch,
				ICallBack: t.ICallBack,
			}
			next = e.Next()
			p.millisecondList.Remove(e)
			continue
		}
		break
	}
}

// AddSecond 添加秒级定时器
//
//	参数:
//		callBackFunc: 回调接口
//		expire: 过期秒数
//	返回值:
//		秒定时器
func (p *defaultTimer) AddSecond(callBackFunc xcontrol.ICallBack, expire int64) *Second {
	t := &Second{
		ICallBack: callBackFunc,
		ISwitch:   xcontrol.NewSwitchButton(true),
		expire:    expire,
	}
	p.secondChan <- t
	return t
}

// DelSecond 删除秒级定时器
// 同 DelMillisecond
func (p *defaultTimer) DelSecond(t *Second) {
	t.reset()
}

// 将秒级定时器,添加到轮转IDX的末尾.之后,移动到合适的位置
//
//		参数:
//			timerSecond: 秒定时器
//			cycleIdx: 轮序号
//	     needMove: 是否需要移动到合适的位置
func (p *defaultTimer) pushBackCycle(timerSecond *Second, cycleIdx int, needMove bool) {
	p.secondSlice[cycleIdx].PushBack(timerSecond)
	if needMove {
		moveLastElementToProperPositionSecond(&p.secondSlice[cycleIdx])
	}
}

// 移动最后一个元素到合适的位置,移动到大于他的元素的前面[实现按照时间排序,加入顺序排序]
// e.g.: 1,2,2,3,4,4,3 => 1,2,2,3,3,4,4 [将最后一个元素移动到4的前面]
func moveLastElementToProperPositionSecond(l *list.List) {
	lastElement := l.Back() // 获取最后一个元素
	target := lastElement.Value.(*Second)
	var element *list.Element
	for element = lastElement.Prev(); element != nil; element = element.Prev() {
		current := element.Value.(*Second)
		if current.expire <= target.expire {
			l.MoveAfter(lastElement, element)
			return
		}
	}
	if element == nil {
		// 如果没有找到比目标小或等于的元素，将目标元素移动到列表的前面
		l.MoveToFront(lastElement)
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
		if t.ISwitch.IsOff() {
			next = e.Next()
			cycle0.Remove(e)
			continue
		}
		if t.expire <= timestamp {
			p.opts.outgoingTimeoutChan <- &EventTimerSecond{
				ISwitch:   t.ISwitch,
				ICallBack: t.ICallBack,
			}
			next = e.Next()
			cycle0.Remove(e)
			continue
		}
		break
	}
	if 0 != cycle0.Len() { // 如果当前的 cycle 中还有元素,则不需要之后的cycle向前移动
		return
	}
	// 更新时间轮,从序号为1的数组开始
	for idx := 1; idx < cycleSize; idx++ {
		c := &p.secondSlice[idx]
		for e := c.Front(); e != nil; e = next {
			t := e.Value.(*Second)
			if t.ISwitch.IsOff() {
				next = e.Next()
				c.Remove(e)
				continue
			}
			if t.expire <= timestamp {
				p.opts.outgoingTimeoutChan <- &EventTimerSecond{
					ISwitch:   t.ISwitch,
					ICallBack: t.ICallBack,
				}
				next = e.Next()
				c.Remove(e)
				continue
			}
			if newIdx := findPrevCycleIdx(t.expire-timestamp, idx); idx != newIdx {
				next = e.Next()
				c.Remove(e)
				p.pushBackCycle(t, newIdx, false)
				continue
			}
			break
		}
		if 0 != c.Len() { // 如果当前的 cycle 中还有元素,则不需要之后的cycle向前移动
			break
		}
	}
}
