package kcp

import (
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xnetcommon "github.com/75912001/xlib/net/common"
	xpacket "github.com/75912001/xlib/packet"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

// ClientOptions contains options to configure a Server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type ClientOptions struct {
	serverAddress    *string // 服务端的地址 e.g.:127.0.0.1:8787
	iOut             xcontrol.IOut
	sendChanCapacity *uint32 // 发送管道容量
	HeaderStrategy   xpacket.IHeaderStrategy
	xnetcommon.ConnOptions
	xnetcommon.KCPOptions
}

// NewClientOptions 新的ClientOptions
func NewClientOptions() *ClientOptions {
	return &ClientOptions{}
}

func (p *ClientOptions) WithAddress(address string) *ClientOptions {
	p.serverAddress = &address
	return p
}
func (p *ClientOptions) WithIOut(iOut xcontrol.IOut) *ClientOptions {
	p.iOut = iOut
	return p
}
func (p *ClientOptions) WithSendChanCapacity(sendChanCapacity uint32) *ClientOptions {
	p.sendChanCapacity = &sendChanCapacity
	return p
}
func (p *ClientOptions) WithHeaderStrategy(strategy xpacket.IHeaderStrategy) *ClientOptions {
	p.HeaderStrategy = strategy
	return p
}

// mergeClientOptions combines the given *ClientOptions into a single *ClientOptions in a last one wins fashion.
// The specified options are merged with the existing options on the Server, with the specified options taking
// precedence.
func mergeClientOptions(opts ...*ClientOptions) *ClientOptions {
	newOptions := NewClientOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.serverAddress != nil {
			newOptions.WithAddress(*opt.serverAddress)
		}
		if opt.iOut != nil {
			newOptions.WithIOut(opt.iOut)
		}
		if opt.sendChanCapacity != nil {
			newOptions.WithSendChanCapacity(*opt.sendChanCapacity)
		}
		if opt.HeaderStrategy != nil {
			newOptions.WithHeaderStrategy(opt.HeaderStrategy)
		}
		newOptions.ConnOptions.Merge(&opt.ConnOptions)
		newOptions.KCPOptions.Merge(&opt.KCPOptions)
	}
	return newOptions
}

// 配置
func clientConfigure(opts *ClientOptions) error {
	if opts.serverAddress == nil {
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
	return nil
}
