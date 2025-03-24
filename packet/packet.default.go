package packet

import (
	xmessage "github.com/75912001/xlib/message"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

// Packet 数据包
type Packet struct {
	Header    *Header           // 包头
	PBMessage proto.Message     // 消息
	IMessage  xmessage.IMessage // 记录该包对应的处理消息
}

// NewPacket 新建数据包
func NewPacket() *Packet {
	return &Packet{}
}

func (p *Packet) WithHeader(header *Header) *Packet {
	p.Header = header
	return p
}

func (p *Packet) WithPBMessage(pb proto.Message) *Packet {
	p.PBMessage = pb
	return p
}

func (p *Packet) WithIMessage(iMessage xmessage.IMessage) *Packet {
	p.IMessage = iMessage
	return p
}

func (p *Packet) Marshal() (data []byte, err error) {
	if p.PBMessage == nil { // 没有消息体
		//return nil, xerror.NotImplemented
		p.Header.Length = HeaderSize
		buf := p.Header.Pack()
		return buf, nil
	} else { // 有消息体
		data, err = proto.Marshal(p.PBMessage)
		if err != nil {
			return nil, errors.WithMessage(err, xruntime.Location())
		}
		p.Header.Length = HeaderSize + uint32(len(data))
		buf := p.Header.Pack()
		copy(buf[HeaderSize:p.Header.Length], data)
		return buf, nil
	}
}
