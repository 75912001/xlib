package xmap

// node 双向链表节点
type node[K comparable, V any] struct {
	key   K
	value V
	prev  *node[K, V]
	next  *node[K, V]
}

// ListOrderedMap 使用双向链表实现的有序映射
type ListOrderedMap[K comparable, V any] struct {
	m  map[K]*node[K, V] // 存储键值对
	lh *node[K, V]       // list head 链表头
	lt *node[K, V]       // list tail 链表尾
}

// NewListOrderedMap 创建新的实例
func NewListOrderedMap[K comparable, V any]() *ListOrderedMap[K, V] {
	return &ListOrderedMap[K, V]{
		m: make(map[K]*node[K, V]),
	}
}

// Add 添加元素
func (p *ListOrderedMap[K, V]) Add(key K, value V) {
	if n, exists := p.m[key]; exists { // 存在
		// 更新
		n.value = value
		return
	}

	// 创建新节点
	n := &node[K, V]{
		key:   key,
		value: value,
	}

	// 添加到链表尾部
	if p.lt == nil {
		p.lh = n
		p.lt = n
	} else {
		n.prev = p.lt
		p.lt.next = n
		p.lt = n
	}

	// 存储到 map
	p.m[key] = n
}

// Del 删除元素
func (p *ListOrderedMap[K, V]) Del(key K) {
	n, exists := p.m[key]
	if !exists { // 不存在
		return
	}

	// 从链表中删除节点
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		p.lh = n.next
	}

	if n.next != nil {
		n.next.prev = n.prev
	} else {
		p.lt = n.prev
	}

	// 从 map 中删除
	delete(p.m, key)
}

// Find 查找元素
func (p *ListOrderedMap[K, V]) Find(key K) (V, bool) {
	if n, exists := p.m[key]; exists { // 存在
		return n.value, true
	}
	var zero V
	return zero, false
}

// Foreach 遍历元素
//
//	f 返回 false 时停止遍历
//	[⚠️] Foreach 中调用增加或移除元素的方法,可能会影响遍历结果.
func (p *ListOrderedMap[K, V]) Foreach(f func(key K, value V) (isContinue bool)) {
	for n := p.lh; n != nil; n = n.next {
		if !f(n.key, n.value) {
			break
		}
	}
}

// ReverseForeach 反向遍历元素
//
//	f 返回 false 时停止遍历
//	[⚠️] ReverseForeach 中调用增加或移除元素的方法,可能会影响遍历结果.
func (p *ListOrderedMap[K, V]) ReverseForeach(f func(key K, value V) (isContinue bool)) {
	for n := p.lt; n != nil; n = n.prev {
		if !f(n.key, n.value) {
			break
		}
	}
}

// Len 获取元素数量
func (p *ListOrderedMap[K, V]) Len() int {
	return len(p.m)
}

// Clear 清空所有元素
func (p *ListOrderedMap[K, V]) Clear() {
	p.lh = nil
	p.lt = nil
	p.m = make(map[K]*node[K, V])
}

// Shrink 优化内存占用
func (p *ListOrderedMap[K, V]) Shrink() {
	// 创建新的 map,容量刚好匹配当前元素数量
	newM := make(map[K]*node[K, V], len(p.m))

	// 复制数据
	for k, v := range p.m {
		newM[k] = v
	}

	// 替换为优化后的数据结构
	p.m = newM
}

// First 获取第一个元素
func (p *ListOrderedMap[K, V]) First() (K, V, bool) {
	if p.lh == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	return p.lh.key, p.lh.value, true
}

// Last 获取最后一个元素
func (p *ListOrderedMap[K, V]) Last() (K, V, bool) {
	if p.lt == nil {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}
	return p.lt.key, p.lt.value, true
}
