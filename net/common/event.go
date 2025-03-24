package common

import (
	xpacket "github.com/75912001/xlib/packet"
)

type IEvent interface {
	Connect(handler IHandler, remote IRemote) error                        // 链接
	Disconnect(handler IHandler, remote IRemote) error                     // 断开链接
	Packet(handler IHandler, remote IRemote, packet xpacket.IPacket) error // 数据包
}

// Connect 事件数据-链接成功
type Connect struct {
	IHandler IHandler
	IRemote  IRemote
}

// Disconnect 事件数据-断开链接
type Disconnect struct {
	IHandler IHandler
	IRemote  IRemote
}

// Packet 事件数据-数据包
type Packet struct {
	IHandler IHandler
	IRemote  IRemote
	IPacket  xpacket.IPacket
}
