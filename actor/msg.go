package actor

import (
	"context"
)

type Msg struct {
	Ctx       context.Context        // 上下文, 可用于取消或设置超时
	Cmd       CMD                    // 消息命令
	Args      []any                  // 消息参数
	asyncChan chan *behaviorResponse // 消息响应通道, 用于同步 [nil 时表示异步,非 nil 时表示同步]
}

// NewMsg 创建新的 Msg
//
//	Ctx: 上下文, Cmd: 命令, Args: 参数
func NewMsg(ctx context.Context, cmd CMD, args ...any) *Msg {
	return &Msg{
		Ctx:  ctx,
		Cmd:  cmd,
		Args: args,
	}
}

// withAsyncChan 设置异步响应通道
func (p *Msg) withAsyncChan(responseChan chan *behaviorResponse) *Msg {
	p.asyncChan = responseChan
	return p
}

// IsSync 是否同步
//
//	如果事件有响应通道, 则认为是同步调用
func (p *Msg) IsSync() bool {
	return p.asyncChan != nil
}
