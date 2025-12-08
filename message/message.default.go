package message

import (
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Message struct {
	xcontrol.ICallBack
	newProtoMessage   func() proto.Message   // 创建新的 proto.Message
	stateSwitch       xcontrol.ISwitchButton // 状态开关-该消息是否启用
	passThroughSwitch xcontrol.ISwitchButton // 透传开关-该消息是否透传
}

func newDefaultMessage(opts *options) *Message {
	return &Message{
		ICallBack:         opts.callback,
		newProtoMessage:   opts.newProtoMessage,
		stateSwitch:       opts.stateSwitch,
		passThroughSwitch: opts.passThroughSwitch,
	}
}

func (p *Message) IsPassThrough() bool {
	return p.passThroughSwitch.IsOn()
}

func (p *Message) Execute() error {
	if p.stateSwitch.IsOff() { // 消息是否禁用
		return xerror.Disable
	}
	return p.ICallBack.Execute()
}

// Marshal 序列化
func (p *Message) Marshal(message proto.Message) (data []byte, err error) {
	data, err = proto.Marshal(message)
	if err != nil {
		return nil, errors.WithMessagef(err, "message marshal %v, %v", message, xruntime.Location())
	}
	return data, nil
}

// Unmarshal 反序列化
//
//	message: 反序列化 得到的 消息
func (p *Message) Unmarshal(data []byte) (message proto.Message, err error) {
	message = p.newProtoMessage()
	err = proto.Unmarshal(data, message)
	if err != nil {
		return nil, errors.WithMessagef(err, "message unmarshal %v, %v", data, xruntime.Location())
	}
	return message, nil
}

// JsonUnmarshal 反序列化
//
//	message: 反序列化 得到的 消息
func (p *Message) JsonUnmarshal(data []byte) (message proto.Message, err error) {
	message = p.newProtoMessage()
	err = protojson.Unmarshal(data, message)
	if err != nil {
		return nil, errors.WithMessagef(err, "message json unmarshal %v, %v", data, xruntime.Location())
	}
	return message, nil
}
