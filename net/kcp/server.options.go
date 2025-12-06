package kcp

import (
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xnetcommon "github.com/75912001/xlib/net/common"
	xpacket "github.com/75912001/xlib/packet"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

// ServerOptions contains options to serverConfigure a Server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type ServerOptions struct {
	listenAddress    *string //监听地址 e.g.:xxx.xxx.xxx.xxx:8899
	iOut             xcontrol.IOut
	sendChanCapacity *uint32 // 发送 channel 大小
	HeaderStrategy   xpacket.IHeaderStrategy
	xnetcommon.ConnOptions
	xnetcommon.KCPOptions
	xnetcommon.PacketLimitOptions
	isActor *bool // 是否是 Actor 模式, 如果是则会使用 Actor 来处理连接的事件 default: false
}

// NewOptions 新的Options
func NewOptions() *ServerOptions {
	return new(ServerOptions)
}

func (p *ServerOptions) WithListenAddress(listenAddress string) *ServerOptions {
	p.listenAddress = &listenAddress
	return p
}

func (p *ServerOptions) WithIOut(iOut xcontrol.IOut) *ServerOptions {
	p.iOut = iOut
	return p
}

func (p *ServerOptions) WithSendChanCapacity(sendChanCapacity uint32) *ServerOptions {
	p.sendChanCapacity = &sendChanCapacity
	return p
}

func (p *ServerOptions) WithHeaderStrategy(strategy xpacket.IHeaderStrategy) *ServerOptions {
	p.HeaderStrategy = strategy
	return p
}

func (p *ServerOptions) WithIsActor(isActor bool) *ServerOptions {
	p.isActor = &isActor
	return p
}

// mergeServerOptions combines the given *ServerOptions into a single *ServerOptions in a last one wins fashion.
// The specified options are merged with the existing options on the Server, with the specified options taking
// precedence.
func mergeServerOptions(opts ...*ServerOptions) *ServerOptions {
	so := NewOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.listenAddress != nil {
			so.WithListenAddress(*opt.listenAddress)
		}
		if opt.iOut != nil {
			so.WithIOut(opt.iOut)
		}
		if opt.sendChanCapacity != nil {
			so.WithSendChanCapacity(*opt.sendChanCapacity)
		}
		if opt.HeaderStrategy != nil {
			so.WithHeaderStrategy(opt.HeaderStrategy)
		}
		so.ConnOptions.Merge(&opt.ConnOptions)
		so.KCPOptions.Merge(&opt.KCPOptions)
		so.PacketLimitOptions.Merge(&opt.PacketLimitOptions)
		if opt.isActor != nil {
			so.WithIsActor(*opt.isActor)
		}
	}
	return so
}

// 配置
func serverConfigure(opts *ServerOptions) error {
	if opts.listenAddress == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.iOut == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.sendChanCapacity == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.HeaderStrategy == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.ConnOptions.Configure() != nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.KCPOptions.Configure() != nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.PacketLimitOptions.Configure() != nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.isActor == nil {
		var isActor = false
		opts.isActor = &isActor
	}
	return nil
}
