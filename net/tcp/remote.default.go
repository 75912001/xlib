package tcp

import (
	"context"
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xnetcommon "github.com/75912001/xlib/net/common"
	xpacket "github.com/75912001/xlib/packet"
	xruntime "github.com/75912001/xlib/runtime"
	xutil "github.com/75912001/xlib/util"
	"github.com/pkg/errors"
	"net"
	"strings"
)

// Remote 远端
type Remote struct {
	Conn *net.TCPConn // 连接
	*xnetcommon.DefaultRemote
}

func NewRemote(Conn *net.TCPConn, sendChan chan any, headerStrategy xpacket.IHeaderStrategy) *Remote {
	remote := &Remote{
		Conn: Conn,
		DefaultRemote: &xnetcommon.DefaultRemote{
			SendChan:       sendChan,
			HeaderStrategy: headerStrategy,
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
	//var err error
	//if err = p.Conn.SetKeepAlive(true); err != nil {
	//	xlog.PrintfErr("SetKeepAlive err:%v", err)
	//}
	//if err = p.Conn.SetKeepAlivePeriod(time.Second * 600); err != nil {
	//	xlog.PrintfErr("SetKeepAlivePeriod err:%v", err)
	//}
	// 禁用 Nagle 算法，提高实时性
	if err := p.Conn.SetNoDelay(true); err != nil {
		xlog.PrintfErr("SetNoDelay err:%v", err)
	}
	if connOptions.ReadBuffer != nil {
		if err := p.Conn.SetReadBuffer(*connOptions.ReadBuffer); err != nil {
			xlog.PrintfErr("WithReadBuffer err:%v", err)
		}
	}
	if connOptions.WriteBuffer != nil {
		if err := p.Conn.SetWriteBuffer(*connOptions.WriteBuffer); err != nil {
			xlog.PrintfErr("WithWriteBuffer err:%v", err)
		}
	}
	ctx := context.Background()
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	p.DefaultRemote.CancelFunc = cancelFunc

	go p.onSend(ctxWithCancel)
	switch p.HeaderStrategy.GetHeaderMode() {
	case xpacket.HeaderModeLengthFirst:
		go p.onRecvLengthFirst(ctxWithCancel, iout, handler)
	case xpacket.HeaderModeMessageIDFirst:
		go p.onRecvMessageIDFirst(ctxWithCancel, iout, handler)
	case xpacket.HeaderModeLengthFirst_WithoutLength:
		go p.onRecvLengthFirst_WithoutLength(ctxWithCancel, iout, handler)
	default:
		xlog.PrintfErr("HeaderMode:%v not support", p.HeaderStrategy.GetHeaderMode())
		panic(xerror.NotSupport)
	}
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
