package pool

import "sync"

// Pool 通用对象池模版
type Pool[T any] struct {
	pool  sync.Pool
	reset func(T)
}

// NewPool 创建一个新的对象池
//
//	creator: 对象创建函数
//	reset: 对象重置函数 nil:不执行
func NewPool[T any](creator func() T, reset func(T)) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return creator()
			},
		},
		reset: reset,
	}
}

// Get 从池中获取一个对象
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put 将对象放回池中
func (p *Pool[T]) Put(x T) {
	if p.reset != nil {
		p.reset(x)
	}
	p.pool.Put(x)
}
