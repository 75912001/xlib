package common

import (
	"fmt"
	xpacket "github.com/75912001/xlib/packet"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

// 将数据 packet 放到 data 中
func PushPacket2Data(data []byte, packet xpacket.IPacket) ([]byte, error) {
	packetData, err := packet.Marshal()
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("packet marshal %v, %v", packet, xruntime.Location()))
	}
	if len(data) == 0 { //当 data len == 0 时候, 直接发送 v.data 数据...
		data = packetData
	} else {
		data = append(data, packetData...)
	}
	return data, nil
}
