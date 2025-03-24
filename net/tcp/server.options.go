package tcp

import (
	xerror "github.com/75912001/xlib/error"
	xcommon "github.com/75912001/xlib/net/common"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

// serverOptions contains options to configure a Server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type serverOptions struct {
	listenAddress    *string            // 127.0.0.1:8787
	eventChan        chan<- interface{} // 待处理的事件
	sendChanCapacity *uint32            // 发送 channel 大小
	connOptions      xcommon.ConnOptions
}

// NewServerOptions 新的ServerOptions
func NewServerOptions() *serverOptions {
	return new(serverOptions)
}

func (p *serverOptions) WithListenAddress(listenAddress string) *serverOptions {
	p.listenAddress = &listenAddress
	return p
}

func (p *serverOptions) WithEventChan(eventChan chan<- interface{}) *serverOptions {
	p.eventChan = eventChan
	return p
}

func (p *serverOptions) WithSendChanCapacity(sendChanCapacity uint32) *serverOptions {
	p.sendChanCapacity = &sendChanCapacity
	return p
}

func (p *serverOptions) WithReadBuffer(readBuffer int) *serverOptions {
	p.connOptions.ReadBuffer = &readBuffer
	return p
}

func (p *serverOptions) WithWriteBuffer(writeBuffer int) *serverOptions {
	p.connOptions.WriteBuffer = &writeBuffer
	return p
}

// mergeServerOptions combines the given *serverOptions into a single *serverOptions in a last one wins fashion.
// The specified options are merged with the existing options on the Server, with the specified options taking
// precedence.
func mergeServerOptions(opts ...*serverOptions) *serverOptions {
	newOptions := NewServerOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.listenAddress != nil {
			newOptions.WithListenAddress(*opt.listenAddress)
		}
		if opt.eventChan != nil {
			newOptions.WithEventChan(opt.eventChan)
		}
		if opt.sendChanCapacity != nil {
			newOptions.WithSendChanCapacity(*opt.sendChanCapacity)
		}
		if opt.connOptions.ReadBuffer != nil {
			newOptions.WithReadBuffer(*opt.connOptions.ReadBuffer)
		}
		if opt.connOptions.WriteBuffer != nil {
			newOptions.WithWriteBuffer(*opt.connOptions.WriteBuffer)
		}
	}
	return newOptions
}

// 配置
func serverConfigure(opts *serverOptions) error {
	if opts.listenAddress == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.eventChan == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.sendChanCapacity == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	return nil
}
