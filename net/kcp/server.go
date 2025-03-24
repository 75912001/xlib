package kcp

import (
	"context"
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
	IEvent   xnetcommon.IEvent
	IHandler xnetcommon.IHandler
	listener *kcp.Listener //监听
	options  *ServerOptions
}

// NewServer 新建服务
func NewServer(handler xnetcommon.IHandler) *Server {
	return &Server{
		IEvent:   nil,
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
	p.IEvent = xnetcommon.NewEvent(p.options.eventChan)

	var err error
	if p.options.fec == nil || !*p.options.fec { //FEC 不启用
		if p.listener, err = kcp.ListenWithOptions(*p.options.listenAddress, p.options.blockCrypt, 0, 0); err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
	} else { //FEC 启用
		if p.listener, err = kcp.ListenWithOptions(*p.options.listenAddress, p.options.blockCrypt, 10, 3); err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
	}
	if p.options.connOptions.WriteBuffer != nil {
		if err := p.listener.SetWriteBuffer(*p.options.connOptions.WriteBuffer); err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
	}
	if p.options.connOptions.ReadBuffer != nil {
		if err := p.listener.SetReadBuffer(*p.options.connOptions.ReadBuffer); err != nil {
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
			if p.options.mtuBytes != nil {
				if !udpSession.SetMtu(*p.options.mtuBytes) {
					xlog.PrintfErr("SetMtu false. mtuBytes:%v", *p.options.mtuBytes)
				}
			}
			if p.options.windowSize != nil {
				udpSession.SetWindowSize(*p.options.windowSize, *p.options.windowSize)
			}
			if p.options.ackNoDelay != nil {
				udpSession.SetACKNoDelay(*p.options.ackNoDelay)
			}
			//Turbo Mode： (1, 10, 2, 1);
			udpSession.SetNoDelay(1, 20, 2, 1)
			go p.handleConn(udpSession)
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

func (p *Server) handleConn(udpSession *kcp.UDPSession) {
	remote := NewRemote(udpSession, make(chan interface{}, *p.options.sendChanCapacity))
	xlog.PrintfInfo("accept from UDPSession:%p, conv:%v, RemoteAddr.Network:%v, RemoteAddr.String:%v, remote:%p",
		udpSession, udpSession.GetConv(), udpSession.RemoteAddr().Network(), udpSession.RemoteAddr().String(), remote)
	if err := p.IEvent.Connect(p.IHandler, remote); err != nil {
		xlog.PrintfErr("event.Connect err:%v", err)
		return
	}
	remote.Start(&p.options.connOptions, p.IEvent, p.IHandler)
}
