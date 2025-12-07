package xmap

import (
	"maps"
	"math/rand"
)

type MapMgr[TKey comparable, TVal any] struct {
	m map[TKey]TVal
}

// NewMapMgr 创建 Mgr 实例
func NewMapMgr[TKey comparable, TVal any]() *MapMgr[TKey, TVal] {
	return &MapMgr[TKey, TVal]{
		m: make(map[TKey]TVal),
	}
}

// Add 添加元素
func (p *MapMgr[TKey, TVal]) Add(key TKey, value TVal) {
	p.m[key] = value
}

// AddIfNotExist 添加元素, 如果不存在则-添加-返回true, 如果存在则-不添加-返回false
func (p *MapMgr[TKey, TVal]) AddIfNotExist(key TKey, value TVal) bool {
	if _, exists := p.m[key]; !exists {
		p.m[key] = value
		return true
	}
	return false
}

// Find 查找元素
func (p *MapMgr[TKey, TVal]) Find(key TKey) (TVal, bool) {
	data, ok := p.m[key]
	return data, ok
}

// IsExist 是否存在
func (p *MapMgr[TKey, TVal]) IsExist(key TKey) bool {
	_, ok := p.m[key]
	return ok
}

// Get 获取元素
func (p *MapMgr[TKey, TVal]) Get(key TKey) TVal {
	return p.m[key]
}

// RadomGet 随机获取
func (p *MapMgr[TKey, TVal]) RadomGet() (TKey, TVal, bool) {
	var zeroK TKey
	var zeroV TVal
	count := len(p.m)
	if count == 0 {
		return zeroK, zeroV, false
	}
	target := rand.Intn(count)
	idx := 0
	for k, v := range p.m {
		if idx == target {
			return k, v, true
		}
		idx++
	}
	return zeroK, zeroV, false // 理论上不会到这里
}

// Del 删除元素
func (p *MapMgr[TKey, TVal]) Del(key ...TKey) {
	for _, k := range key {
		delete(p.m, k)
	}
}

// Len 获取元素数量
func (p *MapMgr[TKey, TVal]) Len() int {
	return len(p.m)
}

// Clear 清空 map
func (p *MapMgr[TKey, TVal]) Clear() {
	p.m = make(map[TKey]TVal)
}

// Foreach 遍历所有元素
//
//	如果 f 返回 false，则停止遍历
//	[⚠️] Foreach 中调用增加或移除元素的方法,可能会影响遍历结果
func (p *MapMgr[TKey, TVal]) Foreach(f func(key TKey, value TVal) (isContinue bool)) {
	for k, v := range p.m {
		if !f(k, v) {
			break
		}
	}
}

// Shrink 优化内存使用
//
//	通过创建新的 map 并复制数据来释放多余内存
func (p *MapMgr[TKey, TVal]) Shrink() {
	newMap := make(map[TKey]TVal, len(p.m))
	maps.Copy(newMap, p.m)
	p.m = newMap
}
