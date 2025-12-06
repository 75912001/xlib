package kcp

import (
	"context"
	xconfig "github.com/75912001/xlib/config"
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xnetcommon "github.com/75912001/xlib/net/common"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"github.com/xtaci/kcp-go/v5"
	"io"
	"runtime/debug"
)

// Server 服务端
type Server struct {
	IHandler xnetcommon.IHandler
	listener *kcp.Listener //监听
	options  *ServerOptions
}

// NewServer 新建服务
func NewServer(handler xnetcommon.IHandler) *Server {
	return &Server{
		IHandler: handler,
		listener: nil,
		options:  nil,
	}
}

// Start 开始
func (p *Server) Start(_ context.Context, opts ...*ServerOptions) error {
	p.options = mergeServerOptions(opts...)
	if err := serverConfigure(p.options); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	var err error
	if p.listener, err = kcp.ListenWithOptions(*p.options.listenAddress,
		p.options.KCPOptions.BlockCrypt, *p.options.KCPOptions.DataShards, *p.options.KCPOptions.ParityShards); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.options.ConnOptions.WriteBuffer != nil {
		if err = p.listener.SetWriteBuffer(*p.options.ConnOptions.WriteBuffer); err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
	}
	if p.options.ConnOptions.ReadBuffer != nil {
		if err = p.listener.SetReadBuffer(*p.options.ConnOptions.ReadBuffer); err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
	}
	go func() {
		defer func() {
			if xruntime.IsRelease() {
				if recoverErr := recover(); recoverErr != nil {
					xlog.PrintErr(xerror.GoroutinePanic, recoverErr, debug.Stack())
				}
			}
			xlog.PrintInfo(xerror.GoroutineDone)
		}()
		for {
			udpSession, acceptErr := p.listener.AcceptKCP()
			if acceptErr != nil {
				xlog.PrintfErr("listener.AcceptKCP err:%v", acceptErr)
				if acceptErr.Error() == io.ErrClosedPipe.Error() {
					xlog.PrintfErr("listener.AcceptKCP io.ErrClosedPipe err:%v", acceptErr)
					return
				}
				continue
			}
			if p.options.KCPOptions.Mtu != nil {
				if !udpSession.SetMtu(*p.options.KCPOptions.Mtu) {
					xlog.PrintfErr("SetMtu false. mtuBytes:%v", *p.options.KCPOptions.Mtu)
				}
			}
			if p.options.KCPOptions.SndWindowSize != nil {
				udpSession.SetWindowSize(*p.options.KCPOptions.SndWindowSize, *p.options.KCPOptions.RcvWindowSize)
			}
			if p.options.KCPOptions.AckNodelay != nil {
				udpSession.SetACKNoDelay(*p.options.KCPOptions.AckNodelay)
			}
			//Turbo Mode： (1, 10, 2, 1);
			//Normal Mode: (1, 20, 2, 1)
			udpSession.SetNoDelay(1, 10, 2, 1)
			go p.handleConn(udpSession, p.options.iOut)
		}
	}()
	return nil
}

// Stop 停止 AcceptTCP
func (p *Server) Stop() {
	if p.listener != nil {
		if err := p.listener.Close(); err != nil {
			xlog.PrintfErr("listener.Close err:%v", err)
		}
		p.listener = nil
	}
}

func (p *Server) handleConn(udpSession *kcp.UDPSession, iOut xcontrol.IOut) {
	remote := NewRemote(udpSession, make(chan interface{}, *p.options.sendChanCapacity), p.options.HeaderStrategy)
	remote.PacketLimit = p.options.NewPacketLimitFunc(p.options.MaxCntPerSec)
	xlog.PrintfInfo("accept from UDPSession:%p, conv:%v, RemoteAddr.Network:%v, RemoteAddr.String:%v, remote:%p",
		udpSession, udpSession.GetConv(), udpSession.RemoteAddr().Network(), udpSession.RemoteAddr().String(), remote)
	if xconfig.GConfigMgr.Base.ProcessingModeIsActor() {
		_ = p.IHandler.OnConnect(remote)
	} else {
		iOut.Send(
			&xnetcommon.Connect{
				IHandler: p.IHandler,
				IRemote:  remote,
			},
		)
	}
	remote.Start(&p.options.ConnOptions, iOut, p.IHandler)
}
