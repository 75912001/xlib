package kcp

import (
	xconfig "github.com/75912001/xlib/config"
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xnetcommon "github.com/75912001/xlib/net/common"
	xpacket "github.com/75912001/xlib/packet"
	xutil "github.com/75912001/xlib/util"
	"github.com/pkg/errors"
	"runtime/debug"
)

// 接收数据-长度在前
func (p *Remote) onRecvLengthFirst(iOut xcontrol.IOut, handler xnetcommon.IHandler) {
	defer func() { //断开链接
		//当 Conn 关闭, 该函数会引发 panic
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
	const minSpace = 1024
	const initSize = 2048
	buf := make([]byte, initSize)
	var readIndex int
	for {
	LoopRead:
		buf = xutil.AdjustBufferSize(buf, readIndex, minSpace, initSize)
		readNum, err := p.UDPSession.Read(buf[readIndex:])
		if nil != err {
			xlog.PrintfInfo("remote:%p err:%v", p, err)
			if p.GetDisconnectReason() == xnetcommon.DisconnectReasonUnknown { // 未设置,就设置为客户端主动断开
				p.SetDisconnectReason(xnetcommon.DisconnectReasonClientShutdown)
			}
			return
		}
		readIndex += readNum
		for {
			lengthSize := p.HeaderStrategy.GetLengthSize()
			if readIndex < int(lengthSize) { // 数据不够,继续读
				// xlog.PrintfErr("remote:%p err:%v", p, xerror.LengthNotEnough)
				goto LoopRead
			}
			packetLength := p.HeaderStrategy.UnpackLength(buf)
			packetAllLength := lengthSize + packetLength
			if err = handler.OnCheckPacketLength(packetAllLength); err != nil {
				if errors.Is(err, xerror.LengthNotEnough) { // 数据不够,继续读
					//xlog.PrintfErr("remote:%p err:%v", p, err)
					goto LoopRead
				} else {
					xlog.PrintfErr("remote:%p err:%v", p, err)
					return
				}
			}
			if readIndex < int(packetAllLength) { // 不够一个完整的数据包,继续读
				goto LoopRead
			}
			//完整的数据包
			if err = handler.OnCheckPacketLimit(p); err != nil {
				xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf, err)
			} else {
				packet, err := handler.OnUnmarshalPacket(p, buf[:packetAllLength])
				if err != nil {
					xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf[:packetAllLength], err)
				} else {
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
			// 移动剩余数据到缓冲区开始
			copy(buf, buf[packetAllLength:readIndex])
			readIndex = readIndex - int(packetAllLength)
			if readIndex == 0 {
				goto LoopRead
			}
		}
	}
}

// 接收数据-消息ID在前
func (p *Remote) onRecvMessageIDFirst(iOut xcontrol.IOut, handler xnetcommon.IHandler) {
	defer func() { //断开链接
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
	const minSpace = 1024
	const initSize = 2048
	buf := make([]byte, initSize)
	var readIndex int
	msgIDSize := p.HeaderStrategy.(xpacket.IHeaderStrategyMessageIDFirst).GetMessageIDSize()
	if msgIDSize != 2 && msgIDSize != 4 { // 只支持 2 字节和 4 字节的消息ID
		xlog.PrintfErr("remote:%p msgIDSize:%v not support", p, msgIDSize)
		return
	}
	for {
	LoopRead:
		buf = xutil.AdjustBufferSize(buf, readIndex, minSpace, initSize)
		readNum, err := p.UDPSession.Read(buf[readIndex:])
		if err != nil {
			xlog.PrintfInfo("remote:%p err:%v", p, err)
			if p.GetDisconnectReason() == xnetcommon.DisconnectReasonUnknown { // 未设置,就设置为客户端主动断开
				p.SetDisconnectReason(xnetcommon.DisconnectReasonClientShutdown)
			}
			return
		}
		readIndex += readNum
		for {
			if readIndex < int(msgIDSize) { // 数据不够,继续读
				//xlog.PrintfErr("remote:%p err:%v", p, xerror.LengthNotEnough)
				goto LoopRead
			}
			var msgID uint32
			switch msgIDSize {
			case 2:
				msgID = uint32(xpacket.GEndian.Uint16(buf[:msgIDSize]))
			case 4:
				msgID = xpacket.GEndian.Uint32(buf[:msgIDSize])
			}
			lengthSize := p.HeaderStrategy.(xpacket.IHeaderStrategyMessageIDFirst).GetLengthSizeByMessageID(msgID)
			if lengthSize != 2 && lengthSize != 4 {
				xlog.PrintfErr("remote:%p lengthSize:%v not support", p, lengthSize)
				return
			}
			if readIndex < int(msgIDSize+lengthSize) { // 数据不够,继续读
				//xlog.PrintfErr("remote:%p err:%v", p, xerror.LengthNotEnough)
				goto LoopRead
			}
			var length uint32
			switch lengthSize {
			case 2:
				length = uint32(xpacket.GEndian.Uint16(buf[msgIDSize:]))
			case 4:
				length = xpacket.GEndian.Uint32(buf[msgIDSize:])
			}
			packetAllLength := msgIDSize + lengthSize + length
			if err = handler.OnCheckPacketLength(packetAllLength); err != nil {
				if errors.Is(err, xerror.LengthNotEnough) { // 数据不够,继续读
					//xlog.PrintfErr("remote:%p err:%v", p, err)
					goto LoopRead
				} else {
					xlog.PrintfErr("remote:%p err:%v", p, err)
					return
				}
			}
			if readIndex < int(packetAllLength) { // 不够一个完整的数据包,继续读
				goto LoopRead
			}
			// 完整的数据包
			if err = handler.OnCheckPacketLimit(p); err != nil {
				xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf, err)
			} else {
				if packet, err := handler.OnUnmarshalPacket(p, buf[:packetAllLength]); err != nil {
					xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf[:packetAllLength], err)
				} else {
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
			// 移动剩余数据到缓冲区开始
			copy(buf, buf[packetAllLength:readIndex])
			readIndex -= int(packetAllLength)
			if readIndex == 0 {
				goto LoopRead
			}
		}
	}
}
