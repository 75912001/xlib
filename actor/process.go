package actor

import (
	"fmt"
	"runtime/debug"
	"time"

	xerror "github.com/75912001/xlib/error"
	xlog "github.com/75912001/xlib/log"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

// process 处理消息的主循环
func (p *Actor[KEY]) process(args ...any) (err error) {
	start := time.Now()
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v stack:%v", r, string(debug.Stack()))
			xlog.PrintfErr("Actor process panic key:%v err:%v", p.key, err)
		}
		p.Statistics.ProcessTime += time.Since(start)
		p.Statistics.Count++
		if err != nil {
			p.Statistics.ErrorCount++
		}
	}()
	var resp any
	msg := args[0]
	switch message := msg.(type) {
	case *Msg: // 处理 Msg
		if isSystemReservedCommand(message.Cmd) { // 系统保留命令
			switch message.Cmd {
			case SystemReservedCommand_Stop:
				resp, err = p.stop(message)
			case SystemReservedCommand_RemoveChild:
				resp, err = p.removeChild(message)
			case SystemReservedCommand_Spawn:
				resp, err = p.spawn(message)
			case SystemReservedCommand_GetChild:
				resp, err = p.getChild(message)
			default:
				err = fmt.Errorf("unknown system reserved command: %v %v", message.Cmd, xruntime.Location())
			}
		} else { // 用户自定义命令
			if p.behavior == nil {
				err = errors.WithMessagef(xerror.NoBehavior, "actor has no behavior. key:%v %v",
					p.key, xruntime.Location())
			} else {
				p.behavior, resp, err = p.behavior(message)
				if err != nil {
					err = errors.WithMessagef(err, "actor process error. %v", xruntime.Location())
				}
			}
		}
		if message.IsSync() { // 同步调用
			message.syncChan <- &behaviorResponse{
				respData: resp,
				err:      err,
			}
		}
		return
	default: // 类型未知
		if p.behavior == nil {
			err = errors.WithMessagef(xerror.NoBehavior, "actor has no behavior. %v", xruntime.Location())
			return
		}
		p.behavior, resp, err = p.behavior(msg)
		if err != nil {
			err = errors.WithMessagef(err, "actor process error. %v", xruntime.Location())
		}
		return
	}
}
