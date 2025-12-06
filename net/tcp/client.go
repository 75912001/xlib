package tcp

import (
	"context"
	xnetcommon "github.com/75912001/xlib/net/common"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"net"
	"time"
)

// Client 客户端
type Client struct {
	IHandler xnetcommon.IHandler
	IRemote  xnetcommon.IRemote
}

func NewClient(handler xnetcommon.IHandler) *Client {
	return &Client{
		IHandler: handler,
		IRemote:  nil,
	}
}

// Connect 连接
func (p *Client) Connect(ctx context.Context, opts ...*ConnectOptions) error {
	opt := mergeConnectOptions(opts...)
	if err := configureConnectOptions(opt); err != nil {
		return errors.WithMessagef(err, "configureConnectOptions:%v %v", opt, xruntime.Location())
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", *opt.serverAddress)
	if nil != err {
		return errors.WithMessagef(err, "ResolveTCPAddr:%v %v", *opt.serverAddress, xruntime.Location())
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if nil != err {
		return errors.WithMessagef(err, "DialTCP:%v %v", tcpAddr, xruntime.Location())
	}

	_ = conn.SetKeepAlive(true)
	_ = conn.SetKeepAlivePeriod(1 * time.Minute)

	remote := NewRemote(conn, make(chan any, *opt.sendChanCapacity), opt.HeaderStrategy)
	p.IRemote = remote
	p.IRemote.Start(&opt.ConnOptions, opt.iOut, p.IHandler)
	return nil
}
