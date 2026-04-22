package id

import (
	"sync"

	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

// IDGenerator 用于在 [min, max] 闭区间内按自增顺序分配 ID。
//
//	用尽后再次 Next 会返回错误，不会回到 min 循环。
//	并发调用 Next 是安全的。
type IDGenerator[K IIDKey] struct {
	mu        sync.Mutex
	current   K
	min       K
	max       K
	exhausted bool
}

// NewIDGenerator 创建发号器。须满足 min < max，否则返回 nil。
func NewIDGenerator[K IIDKey](min, max K) *IDGenerator[K] {
	if !(min < max) {
		return nil
	}
	return &IDGenerator[K]{
		min:     min,
		max:     max,
		current: min,
	}
}

// Next 返回下一个 ID 区间用尽或已处于用尽状态后返回错误。
func (p *IDGenerator[K]) Next() (id K, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.exhausted {
		return 0, errors.WithMessagef(errors.New("id generator exhausted"),
			"min:%v max:%v %v", p.min, p.max, xruntime.Location())
	}

	if p.current < p.min || p.current > p.max {
		p.exhausted = true
		return 0, errors.WithMessagef(errors.New("id out of range"),
			"current:%v min:%v max:%v %v", p.current, p.min, p.max, xruntime.Location())
	}

	id = p.current
	if id == p.max {
		p.exhausted = true
		return id, nil
	}

	p.current++
	return id, nil
}

// 适用于 uint32 和 uint64 类型
type IIDKey interface {
	uint32 | uint64
}
