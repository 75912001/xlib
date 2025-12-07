package actor

import (
	"fmt"
	"time"

	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

// process 处理消息的主循环
func (p *Actor[KEY]) process(args ...any) (err error) {
	start := time.Now()
	defer func() {
		p.Statistics.ProcessTime += time.Since(start)
		p.Statistics.Count++
		if err != nil {
			p.Statistics.ErrorCount++
		}
	}()
	var response any
	msg := args[0]
	switch event := msg.(type) {
	case *BehaviorEvent: // 处理 BehaviorEvent
		if isSystemReservedCommand(event.Cmd) { // 系统保留命令
			switch event.Cmd { // 系统保留命令
			case SystemReservedCommand_Stop:
				response, err = p.handleStop(event)
			case SystemReservedCommand_RemoveChild:
				response, err = p.handleRemoveChild(event)
			case SystemReservedCommand_Spawn:
				response, err = p.handleSpawn(event)
			case SystemReservedCommand_GetChild:
				response, err = p.handleGetChild(event)
			default:
				err = fmt.Errorf("unknown system reserved command: %v %v", event.Cmd, xruntime.Location())
			}
		} else { // 用户自定义命令
			if p.behavior == nil {
				err = errors.WithMessagef(xerror.NoBehavior, "actor has no behavior. key:%v %v",
					p.key, xruntime.Location())
			} else {
				p.behavior, response, err = p.behavior(event)
				if err != nil {
					err = errors.WithMessagef(err, "actor process error. %v", xruntime.Location())
				}
			}
		}
		if event.IsSync() { // 同步调用
			event.responseChan <- &behaviorResponse{
				data: response,
				err:  err,
			}
		}
		return
	default: // 类型未知
		if p.behavior == nil {
			err = errors.WithMessagef(xerror.NoBehavior, "actor has no behavior. %v", xruntime.Location())
			return
		}
		p.behavior, response, err = p.behavior(msg)
		_ = response
		err = errors.WithMessagef(err, "actor process error. %v", xruntime.Location())
		return
	}
}
