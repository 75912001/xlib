package server

import (
	xetcd "github.com/75912001/xlib/etcd"
	xlog "github.com/75912001/xlib/log"
	xnettcp "github.com/75912001/xlib/net/common"
	xruntime "github.com/75912001/xlib/runtime"
	xtimer "github.com/75912001/xlib/timer"
	"time"
)

// Handle todo [重要] issue 在处理 event 时候, 向 eventChan 中插入 事件，注意超出eventChan的上限会阻塞.
func (p *Server) Handle() error {
	//在消费eventChan时可能会往eventChan中写入事件，所以关闭服务时不能close eventChan（造成写入阻塞），通过定时检查eventChan大小来关闭
	for {
		select {
		//case <-GBusChannelCheckChan: // todo menglc
		//	xrlog.GetInstance().Warn("receive GBusChannelCheckChan")
		//	if 0 == len(eventChan) && IsServerStopping() {
		//		xrlog.GetInstance().Warn("server is stopping, stop consume GEventChan with length 0")
		//		return
		//	} else {
		//		xrlog.GetInstance().Warnf("server is stopping, waiting for consume GEventChan with length:%d", len(eventChan))
		//	}
		case value := <-p.BusChannel:
			//TODO [*] 应拿尽拿...
			p.TimeMgr.Update()
			var err error
			switch event := value.(type) {
			case *xnettcp.Connect:
				err = event.IHandler.OnConnect(event.IRemote)
			case *xnettcp.Packet:
				if !event.IRemote.IsConnect() {
					continue
				}
				err = event.IHandler.OnPacket(event.IRemote, event.IPacket)
			case *xnettcp.Disconnect:
				err = event.IHandler.OnDisconnect(event.IRemote)
				if !event.IRemote.IsConnect() {
					continue
				}
				event.IRemote.Stop()
			case *xtimer.EventTimerSecond:
				if event.ISwitch.IsOff() {
					continue
				}
				_ = event.ICallBack.Execute()
			case *xtimer.EventTimerMillisecond:
				if event.ISwitch.IsOff() {
					continue
				}
				_ = event.ICallBack.Execute()
				//kcp server
			//case *xrkcp.EventConnect:
			//	err = event.Remote.Server.GetOnEvent().OnConn(event.Remote)
			//case *xrkcp.EventDisconnect:
			//	err = event.Remote.Server.GetOnEvent().OnDisconnect(event.Remote)
			//case *xrkcp.Packet:
			//	if !event.Remote.IsConn() {
			//		continue
			//	}
			//	err = event.Remote.Server.GetOnEvent().OnPacket(event)
			case *xetcd.Event:
				_ = event.ICallBack.Execute()
			//case *mq_nats.Packet:
			//	err = onNatsFunc(event)
			default:
				xlog.PrintfErr("non-existent event:%value %value", value, event)
			}
			if err != nil {
				p.Log.Errorf("Handle event:%v error:%value", value, err)
			}

			if xruntime.IsDebug() {
				dt := time.Now().Sub(p.TimeMgr.NowTime()).Milliseconds()
				if dt > 50 {
					xlog.PrintfErr("cost time50: %value Millisecond with event type:%T", dt, value)
				} else if dt > 20 {
					xlog.PrintfErr("cost time20: %value Millisecond with event type:%T", dt, value)
				} else if dt > 10 {
					xlog.PrintfErr("cost time10: %value Millisecond with event type:%T", dt, value)
				}
			}
		}
	}
}
