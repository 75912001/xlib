package packet

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
