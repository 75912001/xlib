package pool

// 从 xxl 大小的数据开始才有较好的收益. 比 unuse.byte.go 中的性能提高接近三个数量级
// 其他情况基本一致.

import (
	"sync"
)

const (
	// 内存对齐大小
	alignSize = uint32(8)
)

// alignUp 将大小向上对齐到8字节
func alignUp(size uint32) uint32 {
	return (size + alignSize - 1) & ^(alignSize - 1)

}

// bytePoolSize 字节-内存池的大小级别
type bytePoolSize struct {
	minSize uint32 // 最小内存大小,闭区间
	maxSize uint32 // 最大内存大小,闭区间
	growth  uint32 // 步长
	name    string // 池的名称，便于调试和监控
}

// 定义标准内存池大小
const (
	_256B  = 256        // 256B
	_1KB   = 1 << 10    // 1KB
	_4KB   = 4 << 10    // 4KB
	_16KB  = 16 << 10   // 16KB
	_64KB  = 64 << 10   // 64KB
	_256KB = 256 << 10  // 256KB
	_1MB   = 1024 << 10 // 1024KB
	_4MB   = 4 << 20    // 4MB
)

// 定义内存池配置
var poolConfigs = []bytePoolSize{
	{ // 内存对齐后(alignSize):32,64,96,128,160,192,224,256
		minSize: 1,     // 1B, 1 字节
		maxSize: _256B, // 256B, 256 字节
		growth:  32,    // 32 字节步长
		name:    "xxs", // 超微型池
	},
	{ // 内存对齐后(alignSize:8):264,392,520,648,776,904,1024
		minSize: _256B + 1, // 256B+1, 257 字节
		maxSize: _1KB,      // 1KB, 1024 字节
		growth:  128,       // 128 字节步长
		name:    "xs",      // 微型池
	},
	{ // 内存对齐后(alignSize:8):1032,1544,2056,2568,3080,3592,4096
		minSize: _1KB + 1, // 1KB+1, 1025 字节
		maxSize: _4KB,     // 4KB, 4096 字节
		growth:  512,      // 512 字节步长
		name:    "s",      // 小型池
	},
	{ // 内存对齐后(alignSize:8):4104,6152,8200,10248,12296,14344,16384
		minSize: _4KB + 1, // 4KB+1, 4097 字节
		maxSize: _16KB,    // 16KB, 16384 字节
		growth:  2 << 10,  // 2048 字节步长
		name:    "m",      // 中型池
	},
	{ // 内存对齐后(alignSize:8):16392,24584,32776,40968,49160,57352,65536
		minSize: _16KB + 1, // 16KB+1, 16385 字节
		maxSize: _64KB,     // 64KB, 65536 字节
		growth:  8 << 10,   // 8KB, 8192 字节步长
		name:    "l",       // 大型池
	},
	{ // 内存对齐后(alignSize:8):65544,98312,131080,163848,196616,229384,262144
		minSize: _64KB + 1, // 64KB+1, 65537 字节
		maxSize: _256KB,    // 256KB, 262144 字节
		growth:  32 << 10,  // 32KB, 32768 字节步长
		name:    "xl",      // 超大型池
	},
	{ // 内存对齐后(alignSize:8):262152,393224,524296,655368,786440,917512,1048576
		minSize: _256KB + 1, // 256KB+1, 262145 字节
		maxSize: _1MB,       // 1MB, 1048576 字节
		growth:  128 << 10,  // 128KB, 131072 字节步长
		name:    "xxl",      // 巨型池
	},
	{ // 内存对齐后(alignSize:8):1048584,1572872,2097160,2621448,3145736,3670024,4194304
		minSize: _1MB + 1,  // 1MB+1, 1048577 字节
		maxSize: _4MB,      // 4MB, 4194304 字节
		growth:  512 << 10, // 512KB, 524288 字节步长
		name:    "xxxl",    // 超大巨型池
	},
}

// bytePoolElement 表示单个内存池
type bytePoolElement struct {
	config bytePoolSize
	pools  []sync.Pool
}

// bytePoolMgr 管理所有内存池
type bytePoolMgr struct {
	pools   []*bytePoolElement
	maxSize uint32
}

var defaultManager = &bytePoolMgr{}

func init() {
	// 初始化所有内存池
	var lastSize uint32
	for _, config := range poolConfigs {
		pool := &bytePoolElement{
			config: config,
		}
		pool.initialize(lastSize)
		defaultManager.pools = append(defaultManager.pools, pool)
		lastSize = config.maxSize

		if defaultManager.maxSize < config.maxSize {
			defaultManager.maxSize = config.maxSize
		}
	}
}

// initialize 初始化内存池
func (p *bytePoolElement) initialize(lastSize uint32) {
	// 计算起始大小
	startSize := lastSize + 1
	if startSize == 1 {
		startSize = p.config.growth
	}

	// 计算池数量（向上取整）
	sizeRange := p.config.maxSize - startSize
	poolCount := (sizeRange+p.config.growth-1)/p.config.growth + 1

	// 预分配池数组
	p.pools = make([]sync.Pool, poolCount)

	// 初始化每个池
	for i := uint32(0); i < poolCount; i++ {
		// 计算当前池的大小
		size := startSize + i*p.config.growth
		if i == poolCount-1 {
			size = p.config.maxSize
		}

		// 对齐大小并创建池
		alignedSize := alignUp(size)
		finalSize := alignedSize // 捕获正确的大小
		p.pools[i] = sync.Pool{
			New: func() interface{} {
				return make([]byte, finalSize)
			},
		}
	}
}

// 二分查找找到合适的池
func searchPool(size uint32) *bytePoolElement {
	low, high := 0, len(defaultManager.pools)-1
	for low <= high {
		mid := low + (high-low)/2 //nolint:all // 二分法,2:从中间取

		if defaultManager.pools[mid].config.minSize <= size && size <= defaultManager.pools[mid].config.maxSize {
			return defaultManager.pools[mid]
		}

		if size < defaultManager.pools[mid].config.minSize {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return nil
}

// GetBytes 获取指定大小的字节切片
//
//	[⚠️] 返回的数据,不可以append.否则会丢失原指针,导致无法正确回收
func GetBytes(size uint32) []byte {
	if size <= 0 {
		return nil
	}

	if size > defaultManager.maxSize {
		return make([]byte, size)
	}

	// 二分查找找到合适的池
	pool := searchPool(size)
	if pool == nil {
		panic("byte pool not found")
	}
	return pool.acquire(size)
}

// acquire 从池中获取指定大小的切片
//
//	[⚠️] 返回的数据,不可以append.否则会丢失原指针,导致无法正确回收
func (p *bytePoolElement) acquire(size uint32) []byte {
	// 根据请求大小选择合适的池
	index := (size - p.config.minSize + p.config.growth - 1) / p.config.growth
	if index >= uint32(len(p.pools)) {
		index = uint32(len(p.pools)) - 1
	}
	// 获取对齐后的内存切片并截断到请求大小
	buf := p.pools[index].Get().([]byte)
	return buf[:size]
}

// PutBytes 归还字节切片到池中
func PutBytes(buf []byte) {
	if buf == nil {
		return
	}

	size := uint32(cap(buf))
	if size > defaultManager.maxSize {
		return
	}

	// 二分查找找到合适的池
	pool := searchPool(size)
	if pool == nil {
		panic("byte pool not found")
	}
	pool.release(buf)
}

// release 释放切片到对应的池中
func (p *bytePoolElement) release(buf []byte) {
	// 使用容量来判断应该放回哪个池
	capacity := uint32(cap(buf))
	index := (capacity - p.config.minSize) / p.config.growth
	if index >= uint32(len(p.pools)) {
		index = uint32(len(p.pools)) - 1
	}
	p.pools[index].Put(buf)
}
