package id

import (
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

// IDGenerator 用于生成ID
//
//	生成的 ID 会在 [min, max] 范围内循环分配
//	[❕]该实现仅做 自增,不循环 分配
type IDGenerator[K IIDKey] struct {
	current K
	min     K
	max     K
}

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

// Next 获取下一个ID
func (p *IDGenerator[K]) Next() (id K, err error) {
	if p.current > p.max {
		err = errors.WithMessagef(errors.New("id out of range"), "id current:%v, max:%v %v",
			p.current, p.max, xruntime.Location())
		return 0, err
	}

	id = p.current
	p.current++
	return id, nil
}

// 适用于 uint32 和 uint64 类型
type IIDKey interface {
	uint32 | uint64
}
