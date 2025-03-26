package packet

import "encoding/binary"

// 字节序
type endianMode uint32

const (
	LittleEndian endianMode = 0 // 小端-模式
	BigEndian    endianMode = 1 // 大端-模式
)

var EndianMode = LittleEndian

func SetEndianMode(mode endianMode) {
	EndianMode = mode
}

// IsBigEndian 是否为大端模式
func IsBigEndian() bool {
	return EndianMode == BigEndian
}

// IsLittleEndian 是否为小端模式
func IsLittleEndian() bool {
	return EndianMode == LittleEndian
}

func UnpackUint16(buf []byte) uint16 {
	if IsBigEndian() {
		return binary.BigEndian.Uint16(buf)
	} else {
		return binary.LittleEndian.Uint16(buf)
	}
}

func UnpackUint32(buf []byte) uint32 {
	if IsBigEndian() {
		return binary.BigEndian.Uint32(buf)
	} else {
		return binary.LittleEndian.Uint32(buf)
	}
}

func PackUint16(buf []byte, value uint16) {
	if IsBigEndian() {
		binary.BigEndian.PutUint16(buf, value)
	} else {
		binary.LittleEndian.PutUint16(buf, value)
	}
}

func PackUint32(buf []byte, value uint32) {
	if IsBigEndian() {
		binary.BigEndian.PutUint32(buf, value)
	} else {
		binary.LittleEndian.PutUint32(buf, value)
	}
}
