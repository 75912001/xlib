package xmap

import (
	"maps"
	"math/rand"
	"sync"
)

type MapMutexMgr[TKey comparable, TVal any] struct {
	m  map[TKey]TVal
	mu sync.RWMutex
}

// NewMapMutexMgr 创建 Mgr 实例
func NewMapMutexMgr[TKey comparable, TVal any]() *MapMutexMgr[TKey, TVal] {
	return &MapMutexMgr[TKey, TVal]{
		m: make(map[TKey]TVal),
	}
}

// Add 添加元素
func (p *MapMutexMgr[TKey, TVal]) Add(key TKey, value TVal) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.m[key] = value
}

// AddIfNotExist 添加元素，如果不存在则添加
func (p *MapMutexMgr[TKey, TVal]) AddIfNotExist(key TKey, value TVal) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, exists := p.m[key]; !exists {
		p.m[key] = value
		return true
	}
	return false
}

// Find 查找元素
func (p *MapMutexMgr[TKey, TVal]) Find(key TKey) (TVal, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	data, ok := p.m[key]
	return data, ok
}

// Get 获取元素
func (p *MapMutexMgr[TKey, TVal]) Get(key TKey) TVal {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.m[key]
}

// 随机获取
func (p *MapMutexMgr[TKey, TVal]) RadomGet() (TKey, TVal, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
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
func (p *MapMutexMgr[TKey, TVal]) Del(key ...TKey) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, k := range key {
		delete(p.m, k)
	}
}

// Len 获取元素数量
func (p *MapMutexMgr[TKey, TVal]) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.m)
}

// Clear 清空 map
func (p *MapMutexMgr[TKey, TVal]) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.m = make(map[TKey]TVal)
}

// Foreach 遍历所有元素
//
//	如果 f 返回 false, 则停止遍历
//	[⚠️] f 中不要再次调用含有互斥锁的方法
//	[⚠️] Foreach 中调用增加或移除元素的方法,可能会影响遍历结果
func (p *MapMutexMgr[TKey, TVal]) Foreach(f func(key TKey, value TVal) (isContinue bool)) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for k, v := range p.m {
		if !f(k, v) {
			break
		}
	}
}

// Shrink 优化内存使用
//
//	通过创建新的 map 并复制数据来释放多余内存
func (p *MapMutexMgr[TKey, TVal]) Shrink() {
	p.mu.Lock()
	defer p.mu.Unlock()
	newMap := make(map[TKey]TVal, len(p.m))
	maps.Copy(newMap, p.m)
	p.m = newMap
}
