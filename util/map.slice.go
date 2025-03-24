package util

// OrderedMap 是一个泛型结构体，它结合了 slice 和 map 来保持元素顺序并提供快速查找
// [NOTE]当从 keys 切片中删除元素时，使用了切片的拼接操作来保持剩余元素的顺序。
// 这个操作在切片较大时可能效率不高，但在大多数情况下是可以接受的。
// 如果需要频繁地删除元素，并且关心性能，需要考虑使用其他数据结构，如链表或双向队列，来维护键的顺序。
type OrderedMap[K comparable, V any] struct {
	keys []K
	data map[K]V
}

// NewOrderedMap 创建一个新的 OrderedMap 实例
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		data: make(map[K]V),
	}
}

// Set 向 OrderedMap 中添加一个键值对
func (p *OrderedMap[K, V]) Set(key K, value V) {
	if _, exists := p.data[key]; !exists {
		p.keys = append(p.keys, key)
	}
	p.data[key] = value
}

// Get 从 OrderedMap 中根据键获取值
func (p *OrderedMap[K, V]) Get(key K) (V, bool) {
	value, exists := p.data[key]
	return value, exists
}

// Delete 从 OrderedMap 中删除一个键值对
func (p *OrderedMap[K, V]) Delete(key K) {
	if index, exists := p.findIndex(key); exists {
		// 从切片中删除键（保持顺序的复杂性）
		p.keys = append(p.keys[:index], p.keys[index+1:]...)
		// 从 map 中删除键
		delete(p.data, key)
	}
}

// findIndex 在键切片中查找键的索引，如果不存在则返回 -1 和 false
func (p *OrderedMap[K, V]) findIndex(key K) (int, bool) {
	for i, k := range p.keys {
		if k == key {
			return i, true
		}
	}
	return -1, false
}

// Range 迭代 OrderedMap 中的键值对，按插入顺序
func (p *OrderedMap[K, V]) Range(fn func(key K, value V) bool) {
	for _, key := range p.keys {
		if !fn(key, p.data[key]) {
			break
		}
	}
}

// Len 元素个数
func (p *OrderedMap[K, V]) Len() int {
	return len(p.keys)
}
