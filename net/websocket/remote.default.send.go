package websocket

import (
	"context"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xpacket "github.com/75912001/xlib/packet"
	"github.com/gorilla/websocket"
	"runtime/debug"
)

// 处理发送
func (p *Remote) onSend(ctx context.Context) {
	defer func() {
		// 当 Conn 关闭, 该函数会引发 panic
		if err := recover(); err != nil {
			xlog.PrintErr(xerror.GoroutinePanic, p, err, debug.Stack())
		}
		xlog.PrintInfo(xerror.GoroutineDone, p)
	}()
	var err error
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-p.DefaultRemote.SendChan:
			var data []byte // 待发送数据
			data, err = xpacket.AddPacketToData(data, t.(xpacket.IPacket))
			if err != nil {
				xlog.PrintfErr("push2Data err:%v", err)
				continue
			}
			err = p.Conn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				xlog.PrintfErr("WriteMessage err:%v", err)
				continue
			}
		}
	}
}
