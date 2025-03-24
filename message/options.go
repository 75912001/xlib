package message

import (
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type options struct {
	callback          xcontrol.ICallBack     // [required] 消息回调
	newProtoMessage   func() proto.Message   // [required] 创建新的 proto.Message
	stateSwitch       xcontrol.ISwitchButton // [optional] 状态开关-该消息是否启用 [default]:true
	passThroughSwitch xcontrol.ISwitchButton // [optional] 透传开关-该消息是否透传 [default]:false
}

// NewOptions 创建 options
func NewOptions() *options {
	return &options{}
}

func (p *options) WithHandler(callback xcontrol.ICallBack) *options {
	p.callback = callback
	return p
}

func (p *options) WithNewProtoMessage(newProtoMessage func() proto.Message) *options {
	p.newProtoMessage = newProtoMessage
	return p
}

func (p *options) WithStateSwitch(stateSwitch xcontrol.ISwitchButton) *options {
	p.stateSwitch = stateSwitch
	return p
}

func (p *options) WithPassThroughSwitch(passThroughSwitch xcontrol.ISwitchButton) *options {
	p.passThroughSwitch = passThroughSwitch
	return p
}

func merge(opts ...*options) *options {
	so := NewOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.callback != nil {
			so.callback = opt.callback
		}
		if opt.newProtoMessage != nil {
			so.newProtoMessage = opt.newProtoMessage
		}
		if opt.stateSwitch != nil {
			so.stateSwitch = opt.stateSwitch
		}
		if opt.passThroughSwitch != nil {
			so.passThroughSwitch = opt.passThroughSwitch
		}
	}
	return so
}

// 配置
func configure(opts *options) error {
	if opts.stateSwitch == nil {
		opts.stateSwitch = xcontrol.NewSwitchButton(true)
	}
	if opts.passThroughSwitch == nil {
		opts.passThroughSwitch = xcontrol.NewSwitchButton(false)
	}
	if opts.passThroughSwitch.IsOff() { // 非 透传
		if opts.callback == nil { // 没有处理函数
			return errors.WithMessage(xerror.Param, xruntime.Location())
		}
		if opts.newProtoMessage == nil { // 没有创建消息函数
			return errors.WithMessage(xerror.Param, xruntime.Location())
		}
	}
	return nil
}
