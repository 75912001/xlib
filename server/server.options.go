package server

import (
	xerror "github.com/75912001/xlib/error"
	xetcd "github.com/75912001/xlib/etcd"
	xlog "github.com/75912001/xlib/log"
	xnetcommon "github.com/75912001/xlib/net/common"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

// ServerOptions contains options to configure a Server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type ServerOptions struct {
	TCPHandler      xnetcommon.IHandler
	KCPHandler      xnetcommon.IHandler
	LogCallbackFunc xlog.CallBackFunc
	ETCDCallbackFun xetcd.CallbackFun
}

// NewServerOptions 新的ServerOptions
func NewServerOptions() *ServerOptions {
	return new(ServerOptions)
}

func (p *ServerOptions) WithTCPHandler(handler xnetcommon.IHandler) *ServerOptions {
	p.TCPHandler = handler
	return p
}

func (p *ServerOptions) WithKCPHandler(handler xnetcommon.IHandler) *ServerOptions {
	p.KCPHandler = handler
	return p
}

func (p *ServerOptions) WithLogCallbackFunc(callbackFunc xlog.CallBackFunc) *ServerOptions {
	p.LogCallbackFunc = callbackFunc
	return p
}

func (p *ServerOptions) WithETCDCallbackFun(callbackFun xetcd.CallbackFun) *ServerOptions {
	p.ETCDCallbackFun = callbackFun
	return p
}

// mergeServerOptions combines the given *ServerOptions into a single *ServerOptions in a last one wins fashion.
// The specified options are merged with the existing options on the Server, with the specified options taking
// precedence.
func mergeServerOptions(opts ...*ServerOptions) *ServerOptions {
	newOptions := NewServerOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.TCPHandler != nil {
			newOptions.WithTCPHandler(opt.TCPHandler)
		}
		if opt.KCPHandler != nil {
			newOptions.WithKCPHandler(opt.KCPHandler)
		}
		if opt.LogCallbackFunc != nil {
			newOptions.WithLogCallbackFunc(opt.LogCallbackFunc)
		}
		if opt.ETCDCallbackFun != nil {
			newOptions.WithETCDCallbackFun(opt.ETCDCallbackFun)
		}
	}
	return newOptions
}

// 配置
func serverConfigure(opts *ServerOptions) error {
	if opts.TCPHandler == nil && opts.KCPHandler == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.LogCallbackFunc == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	return nil
}
