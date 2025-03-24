package packet

import (
	"encoding/binary"
	xutil "github.com/75912001/xlib/util"
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
	xutil.PackUint32(p.Length, data[0:])
	xutil.PackUint32(p.MessageID, data[4:])
	xutil.PackUint32(p.SessionID, data[8:])
	xutil.PackUint32(p.ResultID, data[12:])
	xutil.PackUint64(p.Key, data[16:])
	return data
}

func (p *Header) Unpack(data []byte) {
	p.Length = binary.LittleEndian.Uint32(data[0:4])
	p.MessageID = binary.LittleEndian.Uint32(data[4:8])
	p.SessionID = binary.LittleEndian.Uint32(data[8:12])
	p.ResultID = binary.LittleEndian.Uint32(data[12:16])
	p.Key = binary.LittleEndian.Uint64(data[16:HeaderSize])
}
