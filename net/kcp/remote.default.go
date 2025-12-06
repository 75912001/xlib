package kcp

import (
	"context"
	"fmt"
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xnetcommon "github.com/75912001/xlib/net/common"
	xpacket "github.com/75912001/xlib/packet"
	xruntime "github.com/75912001/xlib/runtime"
	xutil "github.com/75912001/xlib/util"
	"github.com/pkg/errors"
	"github.com/xtaci/kcp-go/v5"
)

// Remote 远端信息
type Remote struct {
	UDPSession *kcp.UDPSession
	*xnetcommon.DefaultRemote
}

func NewRemote(udpSession *kcp.UDPSession, sendChan chan interface{}, headerStrategy xpacket.IHeaderStrategy) *Remote {
	remote := &Remote{
		UDPSession: udpSession,
		DefaultRemote: &xnetcommon.DefaultRemote{
			SendChan:       sendChan,
			HeaderStrategy: headerStrategy,
		},
	}
	return remote
}

// GetIP 获取IP地址
func (p *Remote) GetIP() string {
	return p.UDPSession.RemoteAddr().String()
}

// 开始
func (p *Remote) Start(connOptions *xnetcommon.ConnOptions, iout xcontrol.IOut, handler xnetcommon.IHandler) {
	ctx := context.Background()
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	p.DefaultRemote.CancelFunc = cancelFunc

	go p.onSend(ctxWithCancel)
	switch p.HeaderStrategy.GetHeaderMode() {
	case xpacket.HeaderModeLengthFirst:
		go p.onRecvLengthFirst(iout, handler)
	case xpacket.HeaderModeMessageIDFirst:
		go p.onRecvMessageIDFirst(iout, handler)
	default:
		xlog.PrintfErr("HeaderMode:%v not support", p.HeaderStrategy.GetHeaderMode())
		panic(xerror.NotSupport)
	}
}

// IsConnect 是否连接
func (p *Remote) IsConnect() bool {
	return nil != p.UDPSession
}

// Send 发送数据
//
//	[⚠️]必须在 总线/actor 中调用
//	参数:
//		packet: 未序列化的包. [NOTE]该数据会被引用,使用层不可写
func (p *Remote) Send(packet xpacket.IPacket) error {
	if !p.IsConnect() {
		return errors.WithMessage(xerror.Link, xruntime.Location())
	}
	err := xutil.PushEventWithTimeout(p.DefaultRemote.SendChan, packet, xnetcommon.EventAddTimeoutDuration)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("Send packet, PushEventWithTimeout %v", xruntime.Location()))
	}
	return nil
}

// Stop 停止
func (p *Remote) Stop() {
	if p.IsConnect() {
		if err := p.UDPSession.Close(); err != nil {
			xlog.PrintfErr("udpSession close err:%v", err)
		}
		p.UDPSession = nil
	}
	if p.DefaultRemote.CancelFunc != nil {
		p.DefaultRemote.CancelFunc()
		p.DefaultRemote.CancelFunc = nil
	}
}
