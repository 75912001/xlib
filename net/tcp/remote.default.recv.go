package tcp

import (
	"context"
	xconfig "github.com/75912001/xlib/config"
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xnetcommon "github.com/75912001/xlib/net/common"
	xpacket "github.com/75912001/xlib/packet"
	xpool "github.com/75912001/xlib/pool"
	"io"
	"runtime/debug"
)

// 接收数据-长度在前
func (p *Remote) onRecvLengthFirst(ctx context.Context, iOut xcontrol.IOut, handler xnetcommon.IHandler) {
	defer func() { // 断开链接
		// 当 Conn 关闭, 该函数会引发 panic
		if r := recover(); r != nil {
			xlog.PrintErr(xerror.GoroutinePanic, p, r, debug.Stack())
		}
		if xconfig.GConfigMgr.Base.ProcessingModeIsActor() {
			_ = handler.OnDisconnect(p)
		} else {
			iOut.Send(
				&xnetcommon.Disconnect{
					IHandler: handler,
					IRemote:  p,
				},
			)
		}

		xlog.PrintInfo(xerror.GoroutineDone, p)
	}()
	// 消息总长度
	lengthSize := p.HeaderStrategy.GetLengthSize()
	lengthBuf := make([]byte, lengthSize)
	for {
		if _, err := io.ReadFull(p.Conn, lengthBuf); err != nil {
			if !xerror.IsNetErrClosing(err) {
				xlog.PrintfErr("remote:%p err:%v", p, err)
			}
			if p.GetDisconnectReason() == xnetcommon.DisconnectReasonUnknown { // 未设置,就设置为对端主动断开
				p.SetDisconnectReason(xnetcommon.DisconnectReasonPeerShutdown)
			}
			return
		}
		length := p.HeaderStrategy.UnpackLength(lengthBuf)
		if err := handler.OnCheckPacketLength(lengthSize + length); err != nil {
			xlog.PrintfErr("remote:%p OnCheckPacketLength err:%v", p, err)
			p.SetDisconnectReason(xnetcommon.DisconnectReasonClientLogic)
			return
		}
		buf := xpool.GetBytes(lengthSize + length)
		copy(buf, lengthBuf)
		if _, err := io.ReadFull(p.Conn, buf[lengthSize:]); err != nil {
			xlog.PrintfErr("remote:%p err:%v", p, err)
			xpool.PutBytes(buf)
			p.SetDisconnectReason(xnetcommon.DisconnectReasonClientLogic)
			return
		}
		if err := handler.OnCheckPacketLimit(p); err != nil {
			xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf, err)
			xpool.PutBytes(buf)
			continue
		}
		packet, err := handler.OnUnmarshalPacket(p, buf)
		xpool.PutBytes(buf)
		if err != nil {
			xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf, err)
			continue
		}
		if xconfig.GConfigMgr.Base.ProcessingModeIsActor() {
			_ = handler.OnPacket(p, packet)
		} else {
			iOut.Send(
				&xnetcommon.Packet{
					IHandler: handler,
					IRemote:  p,
					IPacket:  packet,
				},
			)
		}
	}
}

// 接收数据-消息ID在前
func (p *Remote) onRecvMessageIDFirst(ctx context.Context, iOut xcontrol.IOut, handler xnetcommon.IHandler) {
	defer func() { // 断开链接
		// 当 Conn 关闭, 该函数会引发 panic
		if r := recover(); r != nil {
			xlog.PrintErr(xerror.GoroutinePanic, p, r, debug.Stack())
		}
		if xconfig.GConfigMgr.Base.ProcessingModeIsActor() {
			_ = handler.OnDisconnect(p)
		} else {
			iOut.Send(
				&xnetcommon.Disconnect{
					IHandler: handler,
					IRemote:  p,
				},
			)
		}

		xlog.PrintInfo(xerror.GoroutineDone, p)
	}()

	// 消息ID
	msgIDSize := p.HeaderStrategy.(xpacket.IHeaderStrategyMessageIDFirst).GetMessageIDSize()
	if msgIDSize != 2 && msgIDSize != 4 {
		xlog.PrintfErr("remote:%p msgIDSize:%v not support", p, msgIDSize)
		return
	}
	msgIDBuf := make([]byte, msgIDSize)
	var msgID uint32
	for {
		if _, err := io.ReadFull(p.Conn, msgIDBuf); err != nil {
			if !xerror.IsNetErrClosing(err) {
				xlog.PrintfErr("remote:%p err:%v", p, err)
			}
			if p.GetDisconnectReason() == xnetcommon.DisconnectReasonUnknown { // 未设置,就设置为对端主动断开
				p.SetDisconnectReason(xnetcommon.DisconnectReasonPeerShutdown)
			}
			return
		}
		switch msgIDSize {
		case 2:
			msgID = uint32(xpacket.GEndian.Uint16(msgIDBuf))
		case 4:
			msgID = xpacket.GEndian.Uint32(msgIDBuf)
		}
		lengthSize := p.HeaderStrategy.(xpacket.IHeaderStrategyMessageIDFirst).GetLengthSizeByMessageID(msgID)
		if lengthSize != 2 && lengthSize != 4 {
			xlog.PrintfErr("remote:%p lengthSize:%v not support", p, lengthSize)
			return
		}
		lengthBuf := make([]byte, lengthSize)
		if _, err := io.ReadFull(p.Conn, lengthBuf); err != nil {
			if !xerror.IsNetErrClosing(err) {
				xlog.PrintfErr("remote:%p err:%v", p, err)
			}
			if p.GetDisconnectReason() == xnetcommon.DisconnectReasonUnknown { // 未设置,就设置为对端主动断开
				p.SetDisconnectReason(xnetcommon.DisconnectReasonPeerShutdown)
			}
			return
		}
		var length uint32
		switch lengthSize {
		case 2:
			length = uint32(xpacket.GEndian.Uint16(lengthBuf))
		case 4:
			length = xpacket.GEndian.Uint32(lengthBuf)
		}
		if err := handler.OnCheckPacketLength(msgIDSize + lengthSize + length); err != nil {
			xlog.PrintfErr("remote:%p OnCheckPacketLength err:%v", p, err)
			p.SetDisconnectReason(xnetcommon.DisconnectReasonClientLogic)
			return
		}
		buf := xpool.GetBytes(msgIDSize + lengthSize + length)
		copy(buf, msgIDBuf)
		copy(buf[msgIDSize:], lengthBuf)
		if _, err := io.ReadFull(p.Conn, buf[msgIDSize+lengthSize:]); err != nil {
			xlog.PrintfErr("remote:%p err:%v", p, err)
			xpool.PutBytes(buf)
			p.SetDisconnectReason(xnetcommon.DisconnectReasonClientLogic)
			return
		}
		if err := handler.OnCheckPacketLimit(p); err != nil {
			xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf, err)
			xpool.PutBytes(buf)
			continue
		}
		packet, err := handler.OnUnmarshalPacket(p, buf)
		xpool.PutBytes(buf)
		if err != nil {
			xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf, err)
			continue
		}
		if xconfig.GConfigMgr.Base.ProcessingModeIsActor() {
			_ = handler.OnPacket(p, packet)
		} else {
			iOut.Send(
				&xnetcommon.Packet{
					IHandler: handler,
					IRemote:  p,
					IPacket:  packet,
				},
			)
		}

	}
}

// 接收数据-长度在前
func (p *Remote) onRecvLengthFirst_WithoutLength(ctx context.Context, iOut xcontrol.IOut, handler xnetcommon.IHandler) {
	defer func() { // 断开链接
		// 当 Conn 关闭, 该函数会引发 panic
		if r := recover(); r != nil {
			xlog.PrintErr(xerror.GoroutinePanic, p, r, debug.Stack())
		}
		if xconfig.GConfigMgr.Base.ProcessingModeIsActor() {
			_ = handler.OnDisconnect(p)
		} else {
			iOut.Send(
				&xnetcommon.Disconnect{
					IHandler: handler,
					IRemote:  p,
				},
			)
		}

		xlog.PrintInfo(xerror.GoroutineDone, p)
	}()
	// 消息总长度
	lengthSize := p.HeaderStrategy.GetLengthSize()
	lengthBuf := make([]byte, lengthSize)
	for {
		if _, err := io.ReadFull(p.Conn, lengthBuf); err != nil {
			if !xerror.IsNetErrClosing(err) {
				xlog.PrintfErr("remote:%p err:%v", p, err)
			}
			if p.GetDisconnectReason() == xnetcommon.DisconnectReasonUnknown { // 未设置,就设置为对端主动断开
				p.SetDisconnectReason(xnetcommon.DisconnectReasonPeerShutdown)
			}
			return
		}
		length := p.HeaderStrategy.UnpackLength(lengthBuf)
		if err := handler.OnCheckPacketLength(lengthSize + length); err != nil {
			xlog.PrintfErr("remote:%p OnCheckPacketLength err:%v", p, err)
			p.SetDisconnectReason(xnetcommon.DisconnectReasonClientLogic)
			return
		}
		buf := xpool.GetBytes(length)
		//copy(buf, lengthBuf)
		if _, err := io.ReadFull(p.Conn, buf); err != nil {
			xlog.PrintfErr("remote:%p err:%v", p, err)
			xpool.PutBytes(buf)
			p.SetDisconnectReason(xnetcommon.DisconnectReasonClientLogic)
			return
		}
		if err := handler.OnCheckPacketLimit(p); err != nil {
			xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf, err)
			xpool.PutBytes(buf)
			continue
		}
		packet, err := handler.OnUnmarshalPacket(p, buf)
		xpool.PutBytes(buf)
		if err != nil {
			xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf, err)
			continue
		}
		if xconfig.GConfigMgr.Base.ProcessingModeIsActor() {
			_ = handler.OnPacket(p, packet)
		} else {
			iOut.Send(
				&xnetcommon.Packet{
					IHandler: handler,
					IRemote:  p,
					IPacket:  packet,
				},
			)
		}
	}
}
