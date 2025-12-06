package websocket

import (
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xnetcommon "github.com/75912001/xlib/net/common"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

type ConnectOptions struct {
	serverAddress    *string // 服务端的地址 e.g.:127.0.0.1:8787
	iOut             xcontrol.IOut
	sendChanCapacity *uint32 // 发送管道容量
	xnetcommon.ConnOptions
}

func NewConnectOptions() *ConnectOptions {
	return &ConnectOptions{}
}

func (p *ConnectOptions) WithAddress(address string) *ConnectOptions {
	p.serverAddress = &address
	return p
}

func (p *ConnectOptions) WithIOut(iOut xcontrol.IOut) *ConnectOptions {
	p.iOut = iOut
	return p
}

func (p *ConnectOptions) WithSendChanCapacity(sendChanCapacity uint32) *ConnectOptions {
	p.sendChanCapacity = &sendChanCapacity
	return p
}

func mergeConnectOptions(opts ...*ConnectOptions) *ConnectOptions {
	newOptions := NewConnectOptions()
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
		newOptions.ConnOptions.Merge(&opt.ConnOptions)
	}
	return newOptions
}

// 配置
func configureConnectOptions(opts *ConnectOptions) error {
	if opts.serverAddress == nil {
		return errors.WithMessagef(xerror.Param, "serverAddress is nil. %v", xruntime.Location())
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
	return nil
}
