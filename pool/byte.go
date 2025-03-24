package pool

import (
	"sync"
)

type bytePool struct {
	min    int //最小范围值
	max    int //最大范围值
	growth int //内存增长值-步长
	pool   []sync.Pool
}

// 小于2048时，按64步长增长.>2048时则按2048长度增长
var bytePoolList = []*bytePool{
	{min: 1, max: 2048, growth: 64},         // (2048-1+1)/64=32(份)
	{min: 2049, max: 65536, growth: 2048},   // (65536-2049 + 1)/2048=31(份)
	{min: 65537, max: 262144, growth: 8192}, // (262144-65537 + 1)/8192=24(份)
}

// 池中 最大范围值
var maxAreaValue int

func init() {
	for i := 0; i < len(bytePoolList); i++ {
		bytePoolList[i].makePool()
		if maxAreaValue < bytePoolList[i].max {
			maxAreaValue = bytePoolList[i].max
		}
	}
}

func (p *bytePool) makePool() {
	poolLen := (p.max - p.min + 1) / p.growth
	p.pool = make([]sync.Pool, poolLen)
	for i := 0; i < poolLen; i++ {
		memSize := (p.min - 1) + (i+1)*p.growth
		p.pool[i] = sync.Pool{
			New: func() interface{} {
				//fmt.Println("sync.Pool New memSize:", memSize)
				return make([]byte, memSize)
			}}
	}
}

func (p *bytePool) makeByteSlice(size int) []byte {
	return p.pool[p.getPosByteSize(size)].Get().([]byte)[:size]
}

func (p *bytePool) getPosByteSize(size int) int {
	return (size - p.min) / p.growth
}

func (p *bytePool) releaseByteSlice(byteBuff []byte) {
	p.pool[p.getPosByteSize(cap(byteBuff))].Put(byteBuff)
}

// MakeByteSlice 分配byte切片
//
//	返回 nil 表示分配失败
//	[NOTE] 返回的数据,不可以append.否则会丢失原指针,导致无法正确回收
func MakeByteSlice(size int) []byte {
	if size <= maxAreaValue {
		for i := 0; i < len(bytePoolList); i++ {
			if size <= bytePoolList[i].max {
				return bytePoolList[i].makeByteSlice(size)
			}
		}
	} else { // 超出定义的上限 不使用池
		return make([]byte, size)
	}
	// 失败
	return nil
}

// ReleaseByteSlice 回收
func ReleaseByteSlice(byteBuff []byte) bool {
	if cap(byteBuff) <= maxAreaValue {
		for i := 0; i < len(bytePoolList); i++ {
			if cap(byteBuff) <= bytePoolList[i].max {
				bytePoolList[i].releaseByteSlice(byteBuff)
				return true
			}
		}
	} else {
		return true
	}
	return false
}
