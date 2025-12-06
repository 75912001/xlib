package timer

import (
	"container/heap"
	"sync/atomic"
)

// 小顶堆

// MillisecondTask 表示一个毫秒级定时任务
//
//	seq 用于保证相同到期时间时先来先出
type MillisecondTask struct {
	expire      int64 // 到期时间戳（毫秒）
	millisecond *Millisecond
	seq         uint64 // 自增序号，保证同到期时间先来先出
}

// MillisecondMinHeap 毫秒-数据-小顶堆
//
//	⚠️只允许通过 heap 包操作，不要直接用 append/sort...
type MillisecondMinHeap []*MillisecondTask

func (p *MillisecondMinHeap) Len() int {
	return len(*p)
}
func (p *MillisecondMinHeap) Less(i, j int) bool {
	if (*p)[i].expire != (*p)[j].expire {
		return (*p)[i].expire < (*p)[j].expire
	}
	return (*p)[i].seq < (*p)[j].seq // 到期时间相同，序号小的先出
}
func (p *MillisecondMinHeap) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *MillisecondMinHeap) Push(x any) {
	*p = append(*p, x.(*MillisecondTask))
}
func (p *MillisecondMinHeap) Pop() any {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[0 : n-1]
	return x
}

// 全局自增序号
var globalSeq uint64 = 0

// NewMilliTask 创建新任务，自动分配序号
func NewMilliTask(expire int64, millisecond *Millisecond) *MillisecondTask {
	return &MillisecondTask{
		expire:      expire,
		millisecond: millisecond,
		seq:         atomic.AddUint64(&globalSeq, 1),
	}
}

// InitMilliTaskHeap 初始化堆
func InitMilliTaskHeap() *MillisecondMinHeap {
	h := &MillisecondMinHeap{}
	heap.Init(h)
	return h
}
