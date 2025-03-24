package tcp

import (
	xerror "github.com/75912001/xlib/error"
	xnetcommon "github.com/75912001/xlib/net/common"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

// ClientOptions contains options to configure a Server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type ClientOptions struct {
	serverAddress    *string            // 服务端的地址 e.g.:127.0.0.1:8787
	eventChan        chan<- interface{} // 外部传递的事件处理管道.连接的事件会放入该管道,以供外部处理
	sendChanCapacity *uint32            // 发送管道容量
	connOptions      *xnetcommon.ConnOptions
}

// NewClientOptions 新的ClientOptions
func NewClientOptions() *ClientOptions {
	return &ClientOptions{
		serverAddress:    nil,
		eventChan:        nil,
		sendChanCapacity: nil,
		connOptions:      xnetcommon.NewConnOptions(),
	}
}

func (p *ClientOptions) WithReadBuffer(readBuffer int) *ClientOptions {
	p.connOptions.ReadBuffer = &readBuffer
	return p
}

func (p *ClientOptions) WithWriteBuffer(writeBuffer int) *ClientOptions {
	p.connOptions.WriteBuffer = &writeBuffer
	return p
}

func (p *ClientOptions) WithAddress(address string) *ClientOptions {
	p.serverAddress = &address
	return p
}

func (p *ClientOptions) WithEventChan(eventChan chan<- interface{}) *ClientOptions {
	p.eventChan = eventChan
	return p
}

func (p *ClientOptions) WithSendChanCapacity(sendChanCapacity uint32) *ClientOptions {
	p.sendChanCapacity = &sendChanCapacity
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
func clientConfigure(opts *ClientOptions) error {
	if opts.serverAddress == nil {
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
