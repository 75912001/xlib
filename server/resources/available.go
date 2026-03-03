package resources

import "sync/atomic"

var GResources = NewResources()

type Resources struct {
	availableLoad atomic.Uint32 // 可用负载
}

// NewResources 新的 Resources
func NewResources() *Resources {
	return &Resources{}
}

// 获取可用负载
func (p *Resources) GetAvailableLoad() uint32 {
	return p.availableLoad.Load()
}

// 设置可用负载
func (p *Resources) SetAvailableLoad(availableLoad uint32) {
	p.availableLoad.Store(availableLoad)
}

// 增加可用负载
func (p *Resources) AddAvailableLoad(delta uint32) {
	p.availableLoad.Add(delta)
}

// 减少可用负载
func (p *Resources) SubAvailableLoad(delta uint32) {
	for {
		old := p.availableLoad.Load()
		var newVal uint32

		if old < delta {
			newVal = 0
		} else {
			newVal = old - delta
		}

		// CAS 操作：比较并交换
		// 如果 p.availableLoad 的当前值仍然等于 old，则将其设置为 newVal 并返回 true。
		// 如果在 Load 之后值被其他协程改变了（即当前值 != old），CAS 会返回 false。
		// 此时循环会继续，重新 Load 最新的值并重试，直到成功为止。
		if p.availableLoad.CompareAndSwap(old, newVal) {
			return
		}
	}
}
