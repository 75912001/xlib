package kcp

import (
	"context"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	netcommon "github.com/75912001/xlib/net/common"
	xpacket "github.com/75912001/xlib/packet"
	xutil "github.com/75912001/xlib/util"
	"io"
	"runtime/debug"
	"time"
)

// 处理发送
func (p *Remote) onSend(ctx context.Context) {
	defer func() {
		//if xruntime.IsRelease() { // 当 Conn 关闭, 该函数会引发 panic
		if err := recover(); err != nil {
			xlog.PrintErr(xerror.GoroutinePanic, p, err, debug.Stack())
		}
		//}
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
	const maxCap = 10240
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-p.DefaultRemote.SendChan:
			data, err = xpacket.AddPacketToData(data, t.(xpacket.IPacket))
			if err != nil {
				xlog.PrintfErr("push2Data err:%v", err)
				continue
			}
			for {
				thisTime = time.Now()
				//超时, 防止, 客户端不 read 数据, 导致主循环无法写入.到p.sendChan中, 阻塞主循环逻辑
				if err := netcommon.UpdateWriteDeadline(p.UDPSession, &lastTime, thisTime, writeTimeOutDuration); err != nil {
					xlog.PrintfErr("updateWriteDeadline remote:%p err:%v", p, err)
				}
				writeCnt, err = p.UDPSession.Write(data)
				if 0 < writeCnt {
					data = xutil.TrimLeftBuffer(data, writeCnt, maxCap)
					if len(data) == 0 {
						break
					}
					xlog.PrintfErr("Conn.Write remote:%p writeCnt:%v remaining:%v", p, writeCnt, len(data))
				}
				for 0 < len(p.DefaultRemote.SendChan) {
					t := <-p.DefaultRemote.SendChan
					data, err = xpacket.AddPacketToData(data, t.(xpacket.IPacket))
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
