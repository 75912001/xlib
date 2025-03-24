package common

import (
	xpacket "github.com/75912001/xlib/packet"
)

type ISend interface {
	Send(packet xpacket.IPacket) error
}
