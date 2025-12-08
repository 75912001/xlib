package server

import (
	xactor "github.com/75912001/xlib/actor"
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xnetcommon "github.com/75912001/xlib/net/common"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"time"
)

// 性能监控阈值
const (
	timeThreshold50ms = 50
	timeThreshold20ms = 20
	timeThreshold10ms = 10
)

// 事件处理结果
type eventResult struct {
	err error // 错误
	dt  int64 // 处理时间,毫秒
}

func (p *Server) behavior(events ...any) (be xactor.Behavior, response any, err error) {
	event := events[0]
	result := p.processEvent(event)
	if result.err != nil {
		xlog.GLog.Error(result.err)
	}

	// 性能监控
	p.monitorPerformance(event, result.dt)
	if result.err != nil {
		return p.behavior, response, errors.WithMessagef(result.err, "processEvent error, with event type:%T", event)
	}
	return p.behavior, response, nil
}

// processEvent 处理事件
func (p *Server) processEvent(value any) eventResult {
	var err error
	beginNow := time.Now()
	switch event := value.(type) {
	case *xnetcommon.Connect:
		err = event.IHandler.OnConnect(event.IRemote)
	case *xnetcommon.Packet:
		if event.IRemote.IsConnect() {
			err = event.IHandler.OnPacket(event.IRemote, event.IPacket)
		}
	case *xnetcommon.Disconnect:
		err = event.IHandler.OnDisconnect(event.IRemote)
		if event.IRemote.IsConnect() {
			event.IRemote.Stop()
		}
	case *xcontrol.Event:
		if event.ISwitch.IsOn() {
			_ = event.ICallBack.Execute()
		}
	default:
		err = xerror.NotSupport
		xlog.GLog.Errorf("non-existent event:%v %v", value, event)
	}
	if err != nil {
		xlog.GLog.Errorf("Handle event:%v error:%v", value, err)
	}
	return eventResult{
		err: err,
		dt:  time.Since(beginNow).Milliseconds(),
	}
}

// monitorPerformance 性能监控
func (p *Server) monitorPerformance(value any, dt int64) {
	if !xruntime.IsDebug() {
		return
	}
	switch {
	case dt > timeThreshold50ms:
		xlog.GLog.Warnf("cost time50ms: %v Millisecond with event type:%T", dt, value)
	case dt > timeThreshold20ms:
		xlog.GLog.Warnf("cost time20ms: %v Millisecond with event type:%T", dt, value)
	case dt > timeThreshold10ms:
		xlog.GLog.Warnf("cost time10ms: %v Millisecond with event type:%T", dt, value)
	}
}
