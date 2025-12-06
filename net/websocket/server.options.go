package websocket

import (
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xnetcommon "github.com/75912001/xlib/net/common"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

type ServerOptions struct {
	pattern          *string // "/projectName/gateway/websocket"
	listenAddress    *string // 127.0.0.1:8787
	iOut             xcontrol.IOut
	sendChanCapacity *uint32 // 发送 channel 大小
	xnetcommon.ConnOptions
	xnetcommon.PacketLimitOptions
}

// NewServerOptions 新的ServerOptions
func NewServerOptions() *ServerOptions {
	return new(ServerOptions)
}

func (p *ServerOptions) WithPattern(pattern string) *ServerOptions {
	p.pattern = &pattern
	return p
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

// mergeServerOptions combines the given *ServerOptions into a single *ServerOptions in a last one wins fashion.
// The specified options are merged with the existing options on the Server, with the specified options taking
// precedence.
func mergeServerOptions(opts ...*ServerOptions) *ServerOptions {
	newOptions := NewServerOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.pattern != nil {
			newOptions.pattern = opt.pattern
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
		newOptions.ConnOptions.Merge(&opt.ConnOptions)
		newOptions.PacketLimitOptions.Merge(&opt.PacketLimitOptions)
	}
	return newOptions
}

// 配置
func configureServerOptions(opts *ServerOptions) error {
	if opts.pattern == nil {
		return errors.WithMessagef(xerror.Param, "pattern is nil. %v", xruntime.Location())
	}
	if opts.listenAddress == nil {
		return errors.WithMessagef(xerror.Param, "listenAddress is nil. %v", xruntime.Location())
	}
	if opts.iOut == nil {
		return errors.WithMessagef(xerror.Param, "iOut is nil. %v", xruntime.Location())
	}
	if opts.sendChanCapacity == nil {
		return errors.WithMessagef(xerror.Param, "sendChanCapacity is nil. %v", xruntime.Location())
	}
	if opts.ConnOptions.Configure() != nil {
		return errors.WithMessagef(xerror.Param, "ConnOptions.Configure() is not nil. %v", xruntime.Location())
	}
	if opts.PacketLimitOptions.Configure() != nil {
		return errors.WithMessagef(xerror.Param, "PacketLimitOptions.Configure() is not nil. %v", xruntime.Location())
	}
	return nil
}
