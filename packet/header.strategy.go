package packet

// 消息头, 解析 策略

//HeaderModeLengthFirst
// length 长度 2:uint16, 4:uint32
// messageID 长度 2:uint16, 4:uint32
// body ...

//HeaderModeMessageIDFirst
// messageID 长度 2:uint16, 4:uint32 ...
// length : 由 messageID 决定 是 uint16 还是 uint32 ...
// body ...

// 消息头策略
type IHeaderStrategy interface {
	GetHeaderMode() HeaderMode         // 获取包头模式
	GetLengthSize() uint32             // 消息体长度 的大小
	UnpackLength(buf []byte) uint32    // 解析消息体长度
	UnpackMessageID(buf []byte) uint32 // 解析消息ID
}

// 消息头策略-消息ID优先
type IHeaderStrategyMessageIDFirst interface {
	GetLengthSizeByMessageID(messageID uint32) uint32 // 根据消息ID,获取消息体长度 的大小
	GetMessageIDSize() uint32                         // 消息ID 的长度
}

// 包头-模式
type HeaderMode int

const (
	HeaderModeLengthFirst               HeaderMode = 0 // 长度在前, 长度的值,包含Length字段自身
	HeaderModeMessageIDFirst            HeaderMode = 1 // 消息ID在前
	HeaderModeLengthFirst_WithoutLength HeaderMode = 2 // 长度在前, 长度的值,不包含Length字段自身
)
