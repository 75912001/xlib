package kcp

import (
	"context"
	xlog "github.com/75912001/xlib/log"
	xnetcommon "github.com/75912001/xlib/net/common"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"github.com/xtaci/kcp-go/v5"
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
func (p *Client) Connect(ctx context.Context, opts ...*ClientOptions) error {
	opt := mergeClientOptions(opts...)
	if err := clientConfigure(opt); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	udpSession, err := kcp.DialWithOptions(*opt.serverAddress,
		opt.KCPOptions.BlockCrypt, *opt.KCPOptions.DataShards, *opt.KCPOptions.ParityShards)
	if nil != err {
		return errors.WithMessage(err, xruntime.Location())
	}
	udpSession.SetWindowSize(*opt.KCPOptions.SndWindowSize, *opt.KCPOptions.RcvWindowSize)
	udpSession.SetNoDelay(*opt.KCPOptions.Nodelay, *opt.KCPOptions.Interval, *opt.KCPOptions.Resend, *opt.KCPOptions.Nc)
	udpSession.SetACKNoDelay(*opt.KCPOptions.AckNodelay)
	udpSession.SetMtu(*opt.KCPOptions.Mtu)
	err = udpSession.SetWriteBuffer(*opt.ConnOptions.WriteBuffer)
	if err != nil {
		xlog.PrintErr("SetWriteBuffer failed", err)
	}
	err = udpSession.SetReadBuffer(*opt.ConnOptions.ReadBuffer)
	if err != nil {
		xlog.PrintErr("SetReadBuffer failed", err)
	}

	remote := NewRemote(udpSession, make(chan interface{}, *opt.sendChanCapacity), opt.HeaderStrategy)
	p.IRemote = remote
	p.IRemote.Start(&opt.ConnOptions, opt.iOut, p.IHandler)
	return nil
}
