package server

import (
	xerror "github.com/75912001/xlib/error"
	xgrpcprotointerceptor "github.com/75912001/xlib/grpc/proto/interceptor"
	xlog "github.com/75912001/xlib/log"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
	"runtime/debug"
)

// Server GRPC服务器
type Server struct {
	GrpcServer *grpc.Server
	listener   net.Listener
}

// NewServer 创建GRPC服务器
func NewServer() *Server {
	return &Server{
		GrpcServer: grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				xgrpcprotointerceptor.ShardKeyServerInterceptor(),
				xgrpcprotointerceptor.TraceServerInterceptor(),
			),
		),
	}
}

// Start 启动服务器
func (p *Server) Start(address string) error {
	var err error
	if p.listener != nil {
		return errors.WithMessage(errors.New("server already started"), xruntime.Location())
	}
	p.listener, err = net.Listen("tcp", address)
	if err != nil {
		return errors.WithMessage(errors.New("failed to listen"), xruntime.Location())
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				xlog.PrintErr(xerror.GoroutinePanic, p, err, debug.Stack())
			}
			xlog.PrintErr(xerror.GoroutineDone, p)
		}()
		if err := p.GrpcServer.Serve(p.listener); err != nil {
			xlog.PrintErr("failed to server:", err, xruntime.Location())
		}
	}()
	return nil
}

// Stop 停止服务器
func (p *Server) Stop() error {
	if p.GrpcServer != nil {
		p.GrpcServer.GracefulStop()
	}
	if p.listener != nil {
		if err := p.listener.Close(); err != nil {
			return errors.WithMessage(errors.New("failed to close listener"), xruntime.Location())
		}
		p.listener = nil
	}
	return nil
}
