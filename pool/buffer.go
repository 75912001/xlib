package pool

import (
	"bytes"
)

var Buffer = NewPool(
	func() *bytes.Buffer {
		return &bytes.Buffer{}
	},
	func(buf *bytes.Buffer) {
		buf.Reset()
	},
)

const sSize = 1024
const mSize = 16 * 1024
const lSize = 64 * 1024

// 定义不同大小的内存池
var (
	// 小型缓冲区池 (<=1KB)
	smallBuffer = NewPool(
		func() *bytes.Buffer {
			return bytes.NewBuffer(make([]byte, 0, sSize))
		},
		func(buf *bytes.Buffer) {
			buf.Reset()
		},
	)

	// 中型缓冲区池 (<=16KB)
	mediumBuffer = NewPool(
		func() *bytes.Buffer {
			return bytes.NewBuffer(make([]byte, 0, mSize))
		},
		func(buf *bytes.Buffer) {
			buf.Reset()
		},
	)

	// 大型缓冲区池 (<=64KB)
	largeBuffer = NewPool(
		func() *bytes.Buffer {
			return bytes.NewBuffer(make([]byte, 0, lSize))
		},
		func(buf *bytes.Buffer) {
			buf.Reset()
		},
	)
)

// GetProperBuffer 根据所需大小返回合适的buffer
//
//	[⚠️] 使用限制:
//		1. 写入数据总量不得超过请求的 size 大小, 否则 buffer 内部扩容会导致 PutBuffer 时分类错误
//		2. 不可对返回的 buffer 进行动态追加写入(超出 size 的部分)
//		3. 适用于已知数据大小的固定写入场景
func GetProperBuffer(size int) *bytes.Buffer {
	switch {
	case size <= sSize:
		return smallBuffer.Get()
	case size <= mSize:
		return mediumBuffer.Get()
	default:
		return largeBuffer.Get()
	}
}

// PutBuffer 将buffer放回适当的池中
//
//	[⚠️] 归还前确保 buffer 未因写入超量而扩容, 否则会被放入错误的池中
func PutBuffer(buf *bytes.Buffer) {
	capacity := buf.Cap()
	switch {
	case capacity <= sSize:
		smallBuffer.Put(buf)
	case capacity <= mSize:
		mediumBuffer.Put(buf)
	default:
		largeBuffer.Put(buf)
	}
}
