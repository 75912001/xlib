package packet

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

// Pack 打包包头, 会分配包头中长度的空间
func (p *Header) Pack() []byte {
	data := make([]byte, p.Length)
	GEndian.PutUint32(data[0:], p.Length)
	GEndian.PutUint32(data[4:], p.MessageID)
	GEndian.PutUint32(data[8:], p.SessionID)
	GEndian.PutUint32(data[12:], p.ResultID)
	GEndian.PutUint64(data[16:], p.Key)
	return data
}

func (p *Header) Unpack(data []byte) {
	p.Length = GEndian.Uint32(data[0:4])
	p.MessageID = GEndian.Uint32(data[4:8])
	p.SessionID = GEndian.Uint32(data[8:12])
	p.ResultID = GEndian.Uint32(data[12:16])
	p.Key = GEndian.Uint64(data[16:HeaderSize])
}
