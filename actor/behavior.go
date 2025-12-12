package actor

// Behavior 定义 actor 的行为函数类型
//
//	每个行为函数接收一个消息,处理它,然后返回下一个行为
type Behavior func(messages ...any) (behavior Behavior, resp any, err error)

// 基础行为实现

// EmptyBehavior 空行为,忽略所有消息
func EmptyBehavior(messages ...any) (behavior Behavior, resp any, err error) {
	return EmptyBehavior, nil, nil
}

type behaviorResponse struct {
	respData any   // 返回数据
	err      error // 错误
}
