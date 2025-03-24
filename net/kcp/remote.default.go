package kcp

import (
	"context"
	"encoding/binary"
	"fmt"
	xconstants "github.com/75912001/xlib/constants"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xnetcommon "github.com/75912001/xlib/net/common"
	xpacket "github.com/75912001/xlib/packet"
	xruntime "github.com/75912001/xlib/runtime"
	xutil "github.com/75912001/xlib/util"
	"github.com/pkg/errors"
	"github.com/xtaci/kcp-go/v5"
	"io"
	"runtime/debug"
	"time"
)

// Remote 远端信息
type Remote struct {
	UDPSession       *kcp.UDPSession
	sendChan         chan interface{} //发送管道
	cancelFunc       context.CancelFunc
	Object           interface{}                 // 保存 应用层数据
	DisconnectReason xnetcommon.DisconnectReason // 断开原因
}

func NewRemote(udpSession *kcp.UDPSession, sendChan chan interface{}) *Remote {
	return &Remote{
		UDPSession: udpSession,
		sendChan:   sendChan,
	}
}

func (p *Remote) GetDisconnectReason() xnetcommon.DisconnectReason {
	return p.DisconnectReason
}

func (p *Remote) SetDisconnectReason(reason xnetcommon.DisconnectReason) {
	p.DisconnectReason = reason
}

// GetIP 获取IP地址
func (p *Remote) GetIP() string {
	return p.UDPSession.RemoteAddr().String()
}

// 开始
func (p *Remote) Start(connOptions *xnetcommon.ConnOptions, event xnetcommon.IEvent, handler xnetcommon.IHandler) {
	ctx := context.Background()
	ctxWithCancel, cancelFunc := context.WithCancel(ctx)
	p.cancelFunc = cancelFunc

	go p.onSend(ctxWithCancel)
	go p.onRecv(event, handler)
}

// IsConnect 是否连接
func (p *Remote) IsConnect() bool {
	return nil != p.UDPSession
}

// Send 发送数据
//
//	[NOTE]必须在 总线 中调用
//	参数:
//		packet: 未序列化的包. [NOTE]该数据会被引用,使用层不可写
func (p *Remote) Send(packet xpacket.IPacket) error {
	if !p.IsConnect() {
		return errors.WithMessage(xerror.Link, xruntime.Location())
	}
	err := xutil.PushEventWithTimeout(p.sendChan, packet, xconstants.BusAddTimeoutDuration)
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
	}
	if p.cancelFunc != nil {
		p.cancelFunc()
	}
	p.cancelFunc = nil
	p.UDPSession = nil
}

// 写超时
//
//	只有超过50%时才更新写截止日期
//	参数:
//		lastTime:上次时间 (可能会更新)
//		thisTime:这次时间
//		writeTimeOutDuration:写超时时长
func (p *Remote) updateWriteDeadline(lastTime *time.Time, thisTime time.Time, writeTimeOutDuration time.Duration) error {
	if (writeTimeOutDuration >> 1) < thisTime.Sub(*lastTime) {
		if err := p.UDPSession.SetWriteDeadline(thisTime.Add(writeTimeOutDuration)); err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
		*lastTime = thisTime
	}
	return nil
}

// 处理发送
func (p *Remote) onSend(ctx context.Context) {
	defer func() {
		if xruntime.IsRelease() {
			// 当 Conn 关闭, 该函数会引发 panic
			if err := recover(); err != nil {
				xlog.PrintErr(xerror.GoroutinePanic, p, err, debug.Stack())
			}
		}
		xlog.PrintInfo(xerror.GoroutineDone, p)
	}()

	//上次时间
	var lastTime time.Time
	//这次时间
	var thisTime time.Time
	//超时, 20包/秒, chan 能容纳 5秒的数据包, 这里设置 为 2 秒超时 (根据 p.sendChan 大小和帧率来设定)
	writeTimeOutDuration := time.Second * 2
	var data []byte // 待发送数据
	var err error
	var writeCnt int
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-p.sendChan:
			data, err = xnetcommon.PushPacket2Data(data, t.(xpacket.IPacket))
			if err != nil {
				xlog.PrintfErr("push2Data err:%v", err)
				continue
			}
			for {
				thisTime = time.Now()
				//超时, 防止, 客户端不 read 数据, 导致主循环无法写入.到p.sendChan中, 阻塞主循环逻辑
				if err := p.updateWriteDeadline(&lastTime, thisTime, writeTimeOutDuration); err != nil {
					xlog.PrintfErr("updateWriteDeadline remote:%p err:%v", p, err)
				}
				writeCnt, err = p.UDPSession.Write(data)
				if 0 < writeCnt {
					data = xutil.RearRangeData(data, writeCnt, 10240)
					if len(data) == 0 {
						break
					} else {
						xlog.PrintfErr("Conn.Write remote:%p writeCnt:%v remaining:%v", p, writeCnt, len(data))
					}
				}
				thisTime = time.Now()
				for 0 < len(p.sendChan) {
					t := <-p.sendChan
					data, err = xnetcommon.PushPacket2Data(data, t.(xpacket.IPacket))
					if err != nil {
						xlog.PrintfErr("push2Data err:%v", err)
						continue
					}
				}
				if nil != err {
					if err.Error() == "timeout" { //网络超时
						xlog.PrintfErr("UDPSession.Write timeOut. remote:%p writeCnt:%v, remaining:%v, err:%v",
							p, writeCnt, len(data), err)
						continue
					}
					if err.Error() == io.ErrClosedPipe.Error() {
						xlog.PrintfErr("UDPSession.Write ErrClosedPipe remote:%p writeCnt:%v, remaining:%v, err:%v",
							p, writeCnt, len(data), err)
					} else {
						xlog.PrintfErr("UDPSession.Write remote:%p writeCnt:%v, remaining:%v, err:%v",
							p, writeCnt, len(data), err)
					}
					break
				}
			}
		}
	}
}

// 接收数据
func (p *Remote) onRecv(event xnetcommon.IEvent, handler xnetcommon.IHandler) {
	defer func() { //断开链接
		if xruntime.IsRelease() {
			// 当 Conn 关闭, 该函数会引发 panic
			if err := recover(); err != nil {
				xlog.PrintErr(xerror.GoroutinePanic, p, err, debug.Stack())
			}
		}
		err := event.Disconnect(handler, p)
		if err != nil {
			xlog.PrintfErr("disconnect err:%v", err)
		}
		xlog.PrintInfo(xerror.GoroutineDone, p)
	}()
	buf := make([]byte, 2048)
	var readIndex int
	var err error
	var readNum int
	for {
	LoopRead:
		readNum, err = p.UDPSession.Read(buf[readIndex:])
		if nil != err {
			xlog.PrintfInfo("remote:%p err:%v", p, err)
			if p.GetDisconnectReason() == xnetcommon.DisconnectReasonUnknown { // 未设置,就设置为客户端主动断开
				p.SetDisconnectReason(xnetcommon.DisconnectReasonClientShutdown)
			} else { // 已设置,就不再设置
			}
			return
		}
		readIndex += readNum
		for {
			packetLength := binary.LittleEndian.Uint32(buf)
			err = handler.OnCheckPacketLength(packetLength)
			if nil != err {
				if errors.Is(err, xerror.LengthNotEnough) { // 数据不够,继续读
					xlog.PrintfErr("remote:%p err:%v", p, err)
					goto LoopRead
				} else {
					xlog.PrintfErr("remote:%p err:%v", p, err)
					return
				}
			}
			//完整的数据包
			data := make([]byte, packetLength)
			copy(data, buf[:packetLength])
			copy(buf, buf[packetLength:readIndex])
			readIndex = readIndex - int(packetLength)

			packet, err := handler.OnUnmarshalPacket(p, data)
			if err != nil {
				xlog.PrintfErr("remote:%p data:%v err:%v", p, data, err)
			} else {
				_ = event.Packet(handler, p, packet)
			}
			if 0 == readIndex {
				goto LoopRead
			}
		}
	}
}
