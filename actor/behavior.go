package actor

import (
	"context"
)

// Behavior 定义 actor 的行为函数类型
//
//	每个行为函数接收一个消息,处理它,然后返回下一个行为
type Behavior func(events ...any) (behavior Behavior, response any, err error)

// 基础行为实现

// EmptyBehavior 空行为,忽略所有消息
func EmptyBehavior(events ...any) (behavior Behavior, response any, err error) {
	return EmptyBehavior, nil, nil
}

type BehaviorEvent struct {
	Ctx          context.Context        // 上下文, 可用于取消或设置超时
	Cmd          uint32                 // 事件命令
	Args         []any                  // 事件参数
	responseChan chan *behaviorResponse // 事件响应通道, 用于同步
}

// NewBehaviorEvent 创建新的 BehaviorEvent
//
//	Ctx: 上下文, Cmd: 事件命令, Args: 事件参数
func NewBehaviorEvent(ctx context.Context, cmd uint32, args ...any) *BehaviorEvent {
	return &BehaviorEvent{
		Ctx:  ctx,
		Cmd:  cmd,
		Args: args,
	}
}

func (p *BehaviorEvent) withResponseChan(responseChan chan *behaviorResponse) *BehaviorEvent {
	p.responseChan = responseChan
	return p
}

// IsSyn 是否同步
//
//	如果事件有响应通道, 则认为是同步调用
func (p *BehaviorEvent) IsSync() bool {
	return p.responseChan != nil
}

type behaviorResponse struct {
	data any   // 返回数据
	err  error // 错误
}
