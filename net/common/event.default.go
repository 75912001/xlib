package common

import (
	xconstants "github.com/75912001/xlib/constants"
	xlog "github.com/75912001/xlib/log"
	xpacket "github.com/75912001/xlib/packet"
	xruntime "github.com/75912001/xlib/runtime"
	xutil "github.com/75912001/xlib/util"
	"github.com/pkg/errors"
)

type Event struct {
	eventChan chan<- interface{} // 待处理的事件
}

func NewEvent(eventChan chan<- interface{}) *Event {
	return &Event{
		eventChan: eventChan,
	}
}

// Connect 连接
func (p *Event) Connect(handler IHandler, remote IRemote) error {
	err := xutil.PushEventWithTimeout(p.eventChan,
		&Connect{
			IHandler: handler,
			IRemote:  remote,
		},
		xconstants.BusAddTimeoutDuration)
	if err != nil {
		xlog.PrintfErr("push Connect failed with eventChan full. remote:%v", remote)
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}

// Disconnect 断开链接
func (p *Event) Disconnect(handler IHandler, remote IRemote) error {
	err := xutil.PushEventWithTimeout(p.eventChan,
		&Disconnect{
			IHandler: handler,
			IRemote:  remote,
		},
		xconstants.BusAddTimeoutDuration)
	if err != nil {
		xlog.PrintfErr("push Disconnect failed with eventChan full. remote:%v", remote)
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}

// Packet 数据包
func (p *Event) Packet(handler IHandler, remote IRemote, packet xpacket.IPacket) error {
	err := xutil.PushEventWithTimeout(p.eventChan,
		&Packet{
			IHandler: handler,
			IRemote:  remote,
			IPacket:  packet,
		},
		xconstants.BusAddTimeoutDuration)
	if err != nil {
		xlog.PrintfErr("push Packet failed with eventChan full. remote:%v packet:%v", remote, packet)
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}
