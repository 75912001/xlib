package websocket

import (
	"context"
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xnetcommon "github.com/75912001/xlib/net/common"
	xpacket "github.com/75912001/xlib/packet"
	xruntime "github.com/75912001/xlib/runtime"
	xutil "github.com/75912001/xlib/util"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"strings"
)

// Remote 远端
type Remote struct {
	Conn *websocket.Conn // 连接
	*xnetcommon.DefaultRemote
}

func NewRemote(Conn *websocket.Conn, sendChan chan any) *Remote {
	remote := &Remote{
		Conn: Conn,
		DefaultRemote: &xnetcommon.DefaultRemote{
			SendChan: sendChan,
		},
	}
	return remote
}

// GetIP 获取IP地址
func (p *Remote) GetIP() string {
	slice := strings.Split(p.Conn.RemoteAddr().String(), ":")
	if len(slice) < 1 {
		return ""
	}
	return slice[0]
}

func (p *Remote) Start(connOptions *xnetcommon.ConnOptions, iout xcontrol.IOut, handler xnetcommon.IHandler) {
	ctx := context.Background()
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	p.DefaultRemote.CancelFunc = cancelFunc

	go p.onSend(ctxWithCancel)
}

// IsConnect 是否连接
func (p *Remote) IsConnect() bool {
	return nil != p.Conn
}

// Send 发送数据
//
//	[⚠️]必须在 总线/actor 中调用
//	参数:
//		packet: 未序列化的包. [NOTE]该数据会被引用,使用层不可写
func (p *Remote) Send(packet xpacket.IPacket) error {
	if !p.IsConnect() {
		xlog.PrintfErr("Send packet, IsConnect is false. %v", xruntime.Location())
		return errors.WithMessagef(xerror.Link, "Send packet, IsConnect is false. %v", xruntime.Location())
	}
	err := xutil.PushEventWithTimeout(p.DefaultRemote.SendChan, packet, xnetcommon.EventAddTimeoutDuration)
	if err != nil {
		xlog.PrintfErr("Send packet, PushEventWithTimeout err:%v", err)
		return errors.WithMessagef(err, "Send packet, PushEventWithTimeout err:%v %v", packet, xruntime.Location())
	}
	return nil
}

func (p *Remote) Stop() {
	if p.IsConnect() {
		err := p.Conn.Close()
		if err != nil {
			xlog.PrintfErr("connect close err:%v", err)
		}
		p.Conn = nil
	}
	if p.DefaultRemote.CancelFunc != nil {
		p.DefaultRemote.CancelFunc()
		p.DefaultRemote.CancelFunc = nil
	}
}
