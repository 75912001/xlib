package tcp

import (
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xnetcommon "github.com/75912001/xlib/net/common"
	xpacket "github.com/75912001/xlib/packet"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

// ServerOptions contains options to configure a Server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type ServerOptions struct {
	listenAddress    *string // 127.0.0.1:8787
	iOut             xcontrol.IOut
	sendChanCapacity *uint32 // 发送 channel 大小
	HeaderStrategy   xpacket.IHeaderStrategy
	xnetcommon.ConnOptions
	xnetcommon.PacketLimitOptions
	isActor *bool // 是否是 Actor 模式, 如果是则会使用 Actor 来处理连接的事件 default: false
}

// NewServerOptions 新的ServerOptions
func NewServerOptions() *ServerOptions {
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
	newOptions := NewServerOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.listenAddress != nil {
			newOptions.WithListenAddress(*opt.listenAddress)
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
		newOptions.PacketLimitOptions.Merge(&opt.PacketLimitOptions)
		if opt.isActor != nil {
			newOptions.WithIsActor(*opt.isActor)
		}
	}
	return newOptions
}

// 配置
func configureServerOptions(opts *ServerOptions) error {
	if opts.listenAddress == nil {
		return errors.WithMessagef(xerror.Param, "listenAddress is nil. %v", xruntime.Location())
	}
	if opts.iOut == nil {
		return errors.WithMessagef(xerror.Param, "iOut is nil. %v", xruntime.Location())
	}
	if opts.sendChanCapacity == nil {
		return errors.WithMessagef(xerror.Param, "sendChanCapacity is nil. %v", xruntime.Location())
	}
	if opts.HeaderStrategy == nil {
		return errors.WithMessagef(xerror.Param, "HeaderStrategy is nil. %v", xruntime.Location())
	}
	if opts.ConnOptions.Configure() != nil {
		return errors.WithMessagef(xerror.Param, "ConnOptions.Configure() is not nil. %v", xruntime.Location())
	}
	if opts.PacketLimitOptions.Configure() != nil {
		return errors.WithMessagef(xerror.Param, "PacketLimitOptions.Configure() is not nil. %v", xruntime.Location())
	}
	if opts.isActor == nil {
		var isActor = false
		opts.isActor = &isActor
	}
	return nil
}
