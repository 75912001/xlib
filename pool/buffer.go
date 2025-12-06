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
//	只适用于固定大小的 buffer
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
