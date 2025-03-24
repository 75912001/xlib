package message

import (
	xcontrol "github.com/75912001/xlib/control"
	"google.golang.org/protobuf/proto"
)

type IMessage interface {
	xcontrol.ICallBack
	Marshal(message proto.Message) (data []byte, err error)
	Unmarshal(data []byte) (message proto.Message, err error)
	JsonUnmarshal(data []byte) (message proto.Message, err error)
}
