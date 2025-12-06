package xmap

import "maps"

import "slices"

// SliceOrderedMap 是一个泛型结构体，它使用 slice (保持元素插入顺序) 和 map 来保持元素
//
//	Deprecated: 使用 ListOrderedMap 代替
//	[❕] 这个操作在切片较大时可能效率不高 但在大多数情况下是可以接受的.
//	如果需要频繁地删除元素 并且关心性能 需要考虑使用其他数据结构, 如链表或双向队列, 来维护键的顺序
type SliceOrderedMap[K comparable, V any] struct {
	s []K     // slice 保持元素插入顺序
	m map[K]V // map 快速查找
}

// NewOrderedMap 创建一个新的实例
func NewOrderedMap[K comparable, V any]() *SliceOrderedMap[K, V] {
	return &SliceOrderedMap[K, V]{
		m: make(map[K]V),
	}
}

// Add 添加一个键值对
func (p *SliceOrderedMap[K, V]) Add(key K, value V) {
	if _, exists := p.m[key]; !exists { // key 不存在
		p.s = append(p.s, key)
	}
	p.m[key] = value
}

// Find 根据键获取值
func (p *SliceOrderedMap[K, V]) Find(key K) (V, bool) {
	value, exists := p.m[key]
	return value, exists
}

// Del 删除一个键值对
func (p *SliceOrderedMap[K, V]) Del(key K) {
	if index, exists := p.findIndex(key); exists { // key 存在
		// 从切片中删除键, 保持插入顺序
		p.s = slices.Delete(p.s, index, index+1)
		// 从 map 中删除键
		delete(p.m, key)
	}
}

// findIndex 在键切片中查找键的索引，如果不存在则返回 -1 和 false
func (p *SliceOrderedMap[K, V]) findIndex(key K) (int, bool) {
	for i, k := range p.s {
		if k == key {
			return i, true
		}
	}
	return -1, false
}

// Foreach 迭代，按插入顺序
//
//	[⚠️] Foreach 中调用增加或移除元素的方法,可能会影响遍历结果
func (p *SliceOrderedMap[K, V]) Foreach(f func(key K, value V) (isContinue bool)) {
	for _, key := range p.s {
		if !f(key, p.m[key]) {
			break
		}
	}
}

// ReverseForeach 反向迭代，按插入顺序
//
//	[⚠️] ReverseForeach 中调用增加或移除元素的方法,可能会影响遍历结果
func (p *SliceOrderedMap[K, V]) ReverseForeach(f func(key K, value V) (isContinue bool)) {
	for i := len(p.s) - 1; i >= 0; i-- {
		if !f(p.s[i], p.m[p.s[i]]) {
			break
		}
	}
}

// Len 元素个数
func (p *SliceOrderedMap[K, V]) Len() int {
	return len(p.s)
}

// Clear 清空所有元素
func (p *SliceOrderedMap[K, V]) Clear() {
	p.s = p.s[:0]
	p.m = make(map[K]V)
}

// Shrink 优化内存占用
func (p *SliceOrderedMap[K, V]) Shrink() {
	// 创建新的切片和map,容量刚好匹配当前元素数量
	newSlice := make([]K, len(p.s))
	newMap := make(map[K]V, len(p.m))
	// 复制数据
	copy(newSlice, p.s)
	maps.Copy(newMap, p.m)
	// 替换为优化后的数据结构
	p.s = newSlice
	p.m = newMap
}
