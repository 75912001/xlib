package websocket

import (
	"context"
	xconfig "github.com/75912001/xlib/config"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xnetcommon "github.com/75912001/xlib/net/common"
	xpacket "github.com/75912001/xlib/packet"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"runtime/debug"
)

// Client 客户端
type Client struct {
	IHandler xnetcommon.IHandler
	IRemote  xnetcommon.IRemote
}

func NewClient(handler xnetcommon.IHandler) *Client {
	return &Client{
		IHandler: handler,
		IRemote:  nil,
	}
}

// Connect 连接
func (p *Client) Connect(ctx context.Context, opts ...*ConnectOptions) error {
	opt := mergeConnectOptions(opts...)
	if err := configureConnectOptions(opt); err != nil {
		return errors.WithMessagef(err, "configureConnectOptions:%v %v", opt, xruntime.Location())
	}

	// 连接到 WebSocket 服务器
	conn, _, err := websocket.DefaultDialer.Dial(*opt.serverAddress, nil)
	if err != nil {
		return errors.WithMessagef(err, "dial %v %v", opt, xruntime.Location())
	}
	//defer conn.Close()

	// 发送消息
	//err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, WebSocket!"))
	//if err != nil {
	//	log.Fatal("发送消息失败:", err)
	//}

	// 读取响应
	//_, message, err := conn.ReadMessage()
	//if err != nil {
	//	log.Fatal("读取消息失败:", err)
	//}
	//log.Printf("收到响应: %s", message)

	remote := NewRemote(conn, make(chan any, *opt.sendChanCapacity))
	p.IRemote = remote
	p.IRemote.Start(&opt.ConnOptions, opt.iOut, p.IHandler)
	go func() {
		defer func() {
			if xruntime.IsRelease() {
				if r := recover(); r != nil {
					xlog.PrintErr(xerror.GoroutinePanic, p, r, debug.Stack())
				}
			}
			if xconfig.GConfigMgr.Base.ProcessingModeIsActor() {
				_ = p.IHandler.OnDisconnect(remote)
			} else {
				opt.iOut.Send(
					&xnetcommon.Disconnect{
						IHandler: p.IHandler,
						IRemote:  remote,
					},
				)
			}
			_ = conn.Close()
		}()

		var messageType int
		var buf []byte
		var packet xpacket.IPacket
		// 处理连接
		for {
			// 读取消息
			messageType, buf, err = conn.ReadMessage()
			if err != nil {
				xlog.PrintfErr("read message fail. err:%v %v", err, debug.Stack())
				if remote.GetDisconnectReason() == xnetcommon.DisconnectReasonUnknown { // 未设置,就设置为客户端主动断开
					remote.SetDisconnectReason(xnetcommon.DisconnectReasonClientShutdown)
				}
				break
			}
			if messageType != websocket.BinaryMessage {
				xlog.PrintfErr("read message fail. err:%v %v", xerror.NotSupport, messageType)
				if remote.GetDisconnectReason() == xnetcommon.DisconnectReasonUnknown { // 未设置,就设置为客户端主动断开
					remote.SetDisconnectReason(xnetcommon.DisconnectReasonClientShutdown)
				}
				break
			}
			if err = p.IHandler.OnCheckPacketLimit(remote); err != nil { // 限流
				xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf, err)
				continue
			}
			packet, err = p.IHandler.OnUnmarshalPacket(remote, buf)
			if err != nil {
				xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf, err)
				continue
			}
			if xconfig.GConfigMgr.Base.ProcessingModeIsActor() {
				_ = p.IHandler.OnPacket(remote, packet)
			} else {
				opt.iOut.Send(
					&xnetcommon.Packet{
						IHandler: p.IHandler,
						IRemote:  remote,
						IPacket:  packet,
					},
				)
			}
		}
	}()
	return nil
}
