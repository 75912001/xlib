package util

// NewMapMgr 创建 Mgr 实例
func NewMapMgr[TKey comparable, TVal interface{}]() *MapMgr[TKey, TVal] {
	return &MapMgr[TKey, TVal]{
		elementMap: make(map[TKey]TVal),
	}
}

type MapMgr[TKey comparable, TVal interface{}] struct {
	elementMap map[TKey]TVal
}

// Add 添加元素
func (p *MapMgr[TKey, TVal]) Add(key TKey, value TVal) {
	p.elementMap[key] = value
}

// Find 查找元素
func (p *MapMgr[TKey, TVal]) Find(key TKey) (TVal, bool) {
	data, ok := p.elementMap[key]
	return data, ok
}

// Del 删除元素
func (p *MapMgr[TKey, TVal]) Del(key TKey) {
	delete(p.elementMap, key)
}
