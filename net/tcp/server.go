package tcp

import (
	"context"
	xconfig "github.com/75912001/xlib/config"
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xnetcommon "github.com/75912001/xlib/net/common"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"net"
	"runtime/debug"
	"time"
)

// Server 服务端
type Server struct {
	IHandler xnetcommon.IHandler
	listener *net.TCPListener //监听
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

// 网络 错误 暂时
func netErrorTemporary(tempDelay time.Duration) (newTempDelay time.Duration) {
	if tempDelay == 0 {
		tempDelay = 5 * time.Millisecond
	} else {
		tempDelay *= 2
	}
	if max := 1 * time.Second; tempDelay > max {
		tempDelay = max
	}
	return tempDelay
}

// Start 运行服务
func (p *Server) Start(_ context.Context, opts ...*ServerOptions) error {
	p.options = mergeServerOptions(opts...)
	if err := configureServerOptions(p.options); err != nil {
		return errors.WithMessagef(err, "configureServerOptions:%v %v", p.options, xruntime.Location())
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", *p.options.listenAddress)
	if nil != err {
		return errors.WithMessagef(err, "ResolveTCPAddr:%v %v", *p.options.listenAddress, xruntime.Location())
	}
	p.listener, err = net.ListenTCP("tcp", tcpAddr)
	if nil != err {
		return errors.WithMessagef(err, "ListenTCP:%v %v", tcpAddr, xruntime.Location())
	}
	go func() {
		defer func() {
			if xruntime.IsRelease() {
				if err := recover(); err != nil {
					xlog.PrintErr(xerror.GoroutinePanic, err, debug.Stack())
				}
			}
			xlog.PrintInfo(xerror.GoroutineDone)
		}()
		var tempDelay time.Duration
		for {
			conn, err := p.listener.AcceptTCP()
			if nil != err {
				if xerror.IsNetErrorTimeout(err) {
					tempDelay = netErrorTemporary(tempDelay)
					xlog.PrintfErr("tempDelay:%v, err:%v", tempDelay, err)
					time.Sleep(tempDelay)
					continue
				}
				xlog.PrintfErr("listen.AcceptTCP, err:%v", err)
				return
			}
			_ = conn.SetKeepAlive(true)
			_ = conn.SetKeepAlivePeriod(1 * time.Minute)
			tempDelay = 0
			go p.handleConn(conn, p.options.iOut)
		}
	}()
	return nil
}

// Stop 停止 AcceptTCP
func (p *Server) Stop() {
	if p.listener != nil {
		err := p.listener.Close()
		if err != nil {
			xlog.PrintfErr("listener close err:%v", err)
		}
		p.listener = nil
	}
}

func (p *Server) handleConn(conn *net.TCPConn, iOut xcontrol.IOut) {
	remote := NewRemote(conn, make(chan any, *p.options.sendChanCapacity), p.options.HeaderStrategy)
	if p.options.NewPacketLimitFunc != nil {
		remote.PacketLimit = p.options.NewPacketLimitFunc(p.options.MaxCntPerSec)
	}
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
