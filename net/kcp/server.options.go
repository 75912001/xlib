package kcp

import (
	xerror "github.com/75912001/xlib/error"
	xnetcommon "github.com/75912001/xlib/net/common"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"github.com/xtaci/kcp-go/v5"
)

// ServerOptions contains options to serverConfigure a Server instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
// TODO 修改源码中 minrto 为10ms
type ServerOptions struct {
	listenAddress    *string            //监听地址 e.g.:xxx.xxx.xxx.xxx:8899
	eventChan        chan<- interface{} // 待处理的事件
	sendChanCapacity *uint32            // 发送 channel 大小
	connOptions      xnetcommon.ConnOptions

	blockCrypt kcp.BlockCrypt //加密,解密
	mtuBytes   *int           //e.g.: 1350
	windowSize *int           //e.g.: 512
	fec        *bool          //是否启用FEC
	ackNoDelay *bool          // send ack immediately for each incoming packet(testing purpose)
}

// NewOptions 新的Options
func NewOptions() *ServerOptions {
	return new(ServerOptions)
}

func (p *ServerOptions) WithListenAddress(listenAddress string) *ServerOptions {
	p.listenAddress = &listenAddress
	return p
}

func (p *ServerOptions) WithEventChan(eventChan chan<- interface{}) *ServerOptions {
	p.eventChan = eventChan
	return p
}

func (p *ServerOptions) WithSendChanCapacity(sendChanCapacity uint32) *ServerOptions {
	p.sendChanCapacity = &sendChanCapacity
	return p
}

func (p *ServerOptions) WithReadBuffer(readBuffer int) *ServerOptions {
	p.connOptions.ReadBuffer = &readBuffer
	return p
}

func (p *ServerOptions) WithWriteBuffer(writeBuffer int) *ServerOptions {
	p.connOptions.WriteBuffer = &writeBuffer
	return p
}

func (p *ServerOptions) WithBlockCrypt(blockCrypt kcp.BlockCrypt) *ServerOptions {
	p.blockCrypt = blockCrypt
	return p
}

func (p *ServerOptions) WithMTUBytes(size int) *ServerOptions {
	p.mtuBytes = &size
	return p
}

func (p *ServerOptions) WithWindowSize(size int) *ServerOptions {
	p.windowSize = &size
	return p
}

func (p *ServerOptions) WithFEC(fec bool) *ServerOptions {
	p.fec = &fec
	return p
}

func (p *ServerOptions) WithAckNoDelay(ackNoDelay bool) *ServerOptions {
	p.ackNoDelay = &ackNoDelay
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
		if opt.eventChan != nil {
			so.WithEventChan(opt.eventChan)
		}
		if opt.sendChanCapacity != nil {
			so.WithSendChanCapacity(*opt.sendChanCapacity)
		}
		if opt.connOptions.ReadBuffer != nil {
			so.WithReadBuffer(*opt.connOptions.ReadBuffer)
		}
		if opt.connOptions.WriteBuffer != nil {
			so.WithWriteBuffer(*opt.connOptions.WriteBuffer)
		}
		if opt.blockCrypt != nil {
			so.WithBlockCrypt(opt.blockCrypt)
		}
		if opt.mtuBytes != nil {
			so.WithMTUBytes(*opt.mtuBytes)
		}
		if opt.windowSize != nil {
			so.WithWindowSize(*opt.windowSize)
		}
		if opt.fec != nil {
			so.WithFEC(*opt.fec)
		}
		if opt.ackNoDelay != nil {
			so.WithAckNoDelay(*opt.ackNoDelay)
		}
	}
	return so
}

// 配置
func serverConfigure(opts *ServerOptions) error {
	if opts.listenAddress == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.eventChan == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.sendChanCapacity == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.blockCrypt == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.mtuBytes == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.windowSize == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.ackNoDelay == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	return nil
}
