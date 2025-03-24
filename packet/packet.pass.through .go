package packet

// PacketPassThrough 透传数据包
type PacketPassThrough struct {
	Header  *Header // 包头
	RawData []byte  // 原始数据(包头+包体)
}

// NewPacketPassThrough 新建-透传数据包
func NewPacketPassThrough() *PacketPassThrough {
	return &PacketPassThrough{}
}

func (p *PacketPassThrough) WithHeader(header *Header) *PacketPassThrough {
	p.Header = header
	return p
}

func (p *PacketPassThrough) Marshal() (data []byte, err error) {
	return p.RawData, nil
}
