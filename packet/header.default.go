package packet

import (
	"encoding/binary"
)

const HeaderSize uint32 = 24
const HeaderLengthFieldSize uint32 = 4 // 包头-总长度-字段 的 大小

// Header 包头
type Header struct {
	Length    uint32 // 总包长度,包含包头＋包体长度
	MessageID uint32 // 命令
	SessionID uint32 // 会话id
	ResultID  uint32 // 结果id
	Key       uint64
}

// NewHeader 新建包头
func NewHeader() *Header {
	return &Header{}
}

func (p *Header) Pack() []byte {
	data := make([]byte, p.Length) // [todo menglc] 这里可以使用内存池,记得回收
	if IsBigEndian() {
		binary.BigEndian.PutUint32(data[0:], p.Length)
		binary.BigEndian.PutUint32(data[4:], p.MessageID)
		binary.BigEndian.PutUint32(data[8:], p.SessionID)
		binary.BigEndian.PutUint32(data[12:], p.ResultID)
		binary.BigEndian.PutUint64(data[16:], p.Key)
	} else {
		binary.LittleEndian.PutUint32(data[0:], p.Length)
		binary.LittleEndian.PutUint32(data[4:], p.MessageID)
		binary.LittleEndian.PutUint32(data[8:], p.SessionID)
		binary.LittleEndian.PutUint32(data[12:], p.ResultID)
		binary.LittleEndian.PutUint64(data[16:], p.Key)
	}
	return data
}

func (p *Header) Unpack(data []byte) {
	if IsBigEndian() {
		p.Length = binary.BigEndian.Uint32(data[0:4])
		p.MessageID = binary.BigEndian.Uint32(data[4:8])
		p.SessionID = binary.BigEndian.Uint32(data[8:12])
		p.ResultID = binary.BigEndian.Uint32(data[12:16])
		p.Key = binary.BigEndian.Uint64(data[16:HeaderSize])
	} else {
		p.Length = binary.LittleEndian.Uint32(data[0:4])
		p.MessageID = binary.LittleEndian.Uint32(data[4:8])
		p.SessionID = binary.LittleEndian.Uint32(data[8:12])
		p.ResultID = binary.LittleEndian.Uint32(data[12:16])
		p.Key = binary.LittleEndian.Uint64(data[16:HeaderSize])
	}
}
