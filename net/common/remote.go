package common

import (
	"context"
	xcontrol "github.com/75912001/xlib/control"
	xpacket "github.com/75912001/xlib/packet"
	"time"
)

type IRemote interface {
	ISend
	IPacketLimit
	IsConnect() bool
	Start(connOptions *ConnOptions, iout xcontrol.IOut, handler IHandler)
	Stop()
	GetIP() string
	GetDisconnectReason() DisconnectReason
	SetDisconnectReason(reason DisconnectReason)
}

type DefaultRemote struct {
	SendChan         chan any // 发送管道
	CancelFunc       context.CancelFunc
	Object           any              // 保存 应用层数据
	DisconnectReason DisconnectReason // 断开原因
	HeaderStrategy   xpacket.IHeaderStrategy
	PacketLimit      IPacketLimit
}

func (p *DefaultRemote) GetDisconnectReason() DisconnectReason {
	return p.DisconnectReason
}

func (p *DefaultRemote) SetDisconnectReason(reason DisconnectReason) {
	p.DisconnectReason = reason
}

func (p *DefaultRemote) IsOverload(cnt uint32, nowTime time.Time) bool {
	if p.PacketLimit == nil {
		return false
	}
	return p.PacketLimit.IsOverload(cnt, nowTime)
}
