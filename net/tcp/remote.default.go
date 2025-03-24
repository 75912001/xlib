package tcp

import (
	"context"
	"encoding/binary"
	xconstants "github.com/75912001/xlib/constants"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xnetcommon "github.com/75912001/xlib/net/common"
	xpacket "github.com/75912001/xlib/packet"
	xpool "github.com/75912001/xlib/pool"
	xruntime "github.com/75912001/xlib/runtime"
	xutil "github.com/75912001/xlib/util"
	"github.com/pkg/errors"
	"io"
	"net"
	"runtime/debug"
	"strings"
	"time"
)

// Remote 远端
type Remote struct {
	Conn             *net.TCPConn     // 连接
	sendChan         chan interface{} // 发送管道
	cancelFunc       context.CancelFunc
	Object           interface{}                 // 保存 应用层数据
	DisconnectReason xnetcommon.DisconnectReason // 断开原因
}

func NewRemote(Conn *net.TCPConn, sendChan chan interface{}) *Remote {
	return &Remote{
		Conn:     Conn,
		sendChan: sendChan,
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
	slice := strings.Split(p.Conn.RemoteAddr().String(), ":")
	if len(slice) < 1 {
		return ""
	}
	return slice[0]
}

func (p *Remote) Start(connOptions *xnetcommon.ConnOptions, event xnetcommon.IEvent, handler xnetcommon.IHandler) {
	//var err error
	//if err = p.Conn.SetKeepAlive(true); err != nil {
	//	xlog.PrintfErr("SetKeepAlive err:%v", err)
	//}
	//if err = p.Conn.SetKeepAlivePeriod(time.Second * 600); err != nil {
	//	xlog.PrintfErr("SetKeepAlivePeriod err:%v", err)
	//}
	//if err := p.Conn.SetNoDelay(true); err != nil {
	//	xlog.PrintfErr("SetNoDelay err:%v", err)
	//}
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
	p.cancelFunc = cancelFunc

	go p.onSend(ctxWithCancel)
	go p.onRecv(event, handler)
}

// IsConnect 是否连接
func (p *Remote) IsConnect() bool {
	return nil != p.Conn
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
		xlog.PrintfErr("Send packet, PushEventWithTimeout err:%v", err)
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}

func (p *Remote) Stop() {
	if p.IsConnect() {
		err := p.Conn.Close()
		if err != nil {
			xlog.PrintfErr("connect close err:%v", err)
		}
	}
	if p.cancelFunc != nil {
		p.cancelFunc()
	}
	p.cancelFunc = nil
	p.Conn = nil
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
		if err := p.Conn.SetWriteDeadline(thisTime.Add(writeTimeOutDuration)); err != nil {
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
	// 上次时间
	var lastTime time.Time
	// 这次时间
	var thisTime time.Time
	// 超时
	writeTimeOutDuration := time.Millisecond * 100
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
				// 超时, 防止, 客户端不 read 数据, 导致主循环无法写入.到p.sendChan中, 阻塞主循环逻辑
				if err := p.updateWriteDeadline(&lastTime, thisTime, writeTimeOutDuration); err != nil {
					xlog.PrintfErr("updateWriteDeadline remote:%p err:%v", p, err)
				}
				writeCnt, err = p.Conn.Write(data)
				if 0 < writeCnt {
					data = xutil.RearRangeData(data, writeCnt, 10240)
					if len(data) == 0 {
						break
					} else {
						xlog.PrintfErr("Conn.Write remote:%p writeCnt:%v remaining:%v", p, writeCnt, len(data))
					}
				}
				for 0 < len(p.sendChan) { // 尽量取出待发送数据
					t := <-p.sendChan
					data, err = xnetcommon.PushPacket2Data(data, t.(xpacket.IPacket))
					if err != nil {
						xlog.PrintfErr("push2Data err:%v", err)
						continue
					}
				}
				if nil != err {
					if xutil.IsNetErrorTimeout(err) { // 网络超时
						xlog.PrintfErr("Conn.Write timeOut. remote:%p writeCnt:%v remaining:%v err:%v",
							p, writeCnt, len(data), err)
						continue
					}
					xlog.PrintfErr("Conn.Write remote:%p writeCnt:%v remaining:%v err:%v", p, writeCnt, len(data), err)
					break
				}
			}
		}
	}
}

// 接收数据
func (p *Remote) onRecv(event xnetcommon.IEvent, handler xnetcommon.IHandler) {
	defer func() { // 断开链接
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
	// 消息总长度
	msgLengthBuf := make([]byte, xpacket.HeaderLengthFieldSize)
	for {
		if _, err := io.ReadFull(p.Conn, msgLengthBuf); err != nil {
			if !xutil.IsNetErrClosing(err) {
				xlog.PrintfErr("remote:%p err:%v", p, err)
			}
			if p.GetDisconnectReason() == xnetcommon.DisconnectReasonUnknown { // 未设置,就设置为客户端主动断开
				p.SetDisconnectReason(xnetcommon.DisconnectReasonClientShutdown)
			} else { // 已设置,就不再设置
			}
			return
		}
		packetLength := binary.LittleEndian.Uint32(msgLengthBuf)
		if err := handler.OnCheckPacketLength(packetLength); err != nil {
			xlog.PrintfErr("remote:%p OnCheckPacketLength err:%v", p, err)
			p.SetDisconnectReason(xnetcommon.DisconnectReasonClientLogic)
			return
		}
		buf := xpool.MakeByteSlice(int(packetLength))
		copy(buf, msgLengthBuf)
		if _, err := io.ReadFull(p.Conn, buf[xpacket.HeaderLengthFieldSize:]); err != nil {
			xlog.PrintfErr("remote:%p err:%v", p, err)
			_ = xpool.ReleaseByteSlice(buf)
			p.SetDisconnectReason(xnetcommon.DisconnectReasonClientLogic)
			return
		}
		if err := handler.OnCheckPacketLimit(p); err != nil {
			xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf, err)
			_ = xpool.ReleaseByteSlice(buf)
			continue
		}
		packet, err := handler.OnUnmarshalPacket(p, buf)
		_ = xpool.ReleaseByteSlice(buf)
		if err != nil {
			xlog.PrintfErr("remote:%p buf:%v err:%v", p, buf, err)
			continue
		}
		buf = nil
		_ = event.Packet(handler, p, packet)
	}
}
