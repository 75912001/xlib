package tcp

import (
	"context"
	xnetcommon "github.com/75912001/xlib/net/common"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"net"
)

// Client 客户端
type Client struct {
	IEvent   xnetcommon.IEvent
	IHandler xnetcommon.IHandler
	IRemote  xnetcommon.IRemote
}

func NewClient(handler xnetcommon.IHandler) *Client {
	return &Client{
		IEvent:   nil,
		IHandler: handler,
		IRemote:  nil,
	}
}

// Connect 连接
//
//	每个连接有 一个 发送协程, 一个 接收协程
func (p *Client) Connect(ctx context.Context, opts ...*ClientOptions) error {
	newOpts := mergeClientOptions(opts...)
	if err := clientConfigure(newOpts); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	p.IEvent = xnetcommon.NewEvent(newOpts.eventChan)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", *newOpts.serverAddress)
	if nil != err {
		return errors.WithMessage(err, xruntime.Location())
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if nil != err {
		return errors.WithMessage(err, xruntime.Location())
	}
	p.IRemote = NewRemote(conn, make(chan interface{}, *newOpts.sendChanCapacity))
	p.IRemote.Start(newOpts.connOptions, p.IEvent, p.IHandler)
	return nil
}
