package server

import (
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xetcd "github.com/75912001/xlib/etcd"
	xnetcommon "github.com/75912001/xlib/net/common"
	xpacket "github.com/75912001/xlib/packet"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

type Options struct {
	TCPHandler       xnetcommon.IHandler
	KCPHandler       xnetcommon.IHandler
	WebsocketHandler xnetcommon.IHandler
	LogCallback      xcontrol.ICallBack
	HeaderStrategy   xpacket.IHeaderStrategy
	Etcd             *xetcd.Options
}

// NewServerOptions 新的ServerOptions
func NewServerOptions() *Options {
	opt := &Options{
		Etcd: xetcd.NewOptions(),
	}
	return opt
}

func (p *Options) WithTCPHandler(handler xnetcommon.IHandler) *Options {
	p.TCPHandler = handler
	return p
}

func (p *Options) WithKCPHandler(handler xnetcommon.IHandler) *Options {
	p.KCPHandler = handler
	return p
}

func (p *Options) WithWebsocketHandler(handler xnetcommon.IHandler) *Options {
	p.WebsocketHandler = handler
	return p
}

func (p *Options) WithLogCallbackFunc(callback xcontrol.ICallBack) *Options {
	p.LogCallback = callback
	return p
}

func (p *Options) WithHeaderStrategy(strategy xpacket.IHeaderStrategy) *Options {
	p.HeaderStrategy = strategy
	return p
}

func (p *Options) WithEtcd(etcd *xetcd.Options) *Options {
	p.Etcd = etcd
	return p
}

func mergeOptions(opts ...*Options) *Options {
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
		if opt.WebsocketHandler != nil {
			newOptions.WithWebsocketHandler(opt.WebsocketHandler)
		}
		if opt.LogCallback != nil {
			newOptions.WithLogCallbackFunc(opt.LogCallback)
		}
		if opt.HeaderStrategy != nil {
			newOptions.WithHeaderStrategy(opt.HeaderStrategy)
		}
		if opt.Etcd != nil {
			newOptions.WithEtcd(opt.Etcd)
		}
	}
	return newOptions
}

// 配置
func configure(opts *Options) error {
	if opts.TCPHandler == nil && opts.KCPHandler == nil && opts.WebsocketHandler == nil {
		return errors.WithMessagef(xerror.Param, "tcpHandler and kcpHandler and websocketHandler are nil. %v", xruntime.Location())
	}
	if opts.LogCallback == nil {
		return errors.WithMessagef(xerror.Param, "logCallback is nil. %v", xruntime.Location())
	}
	if opts.HeaderStrategy == nil {
		return errors.WithMessagef(xerror.Param, "headerStrategy is nil. %v", xruntime.Location())
	}
	return nil
}
