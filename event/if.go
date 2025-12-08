package event

// IEvent 事件接口
type IEvent interface {
	Send(event any) // 发送事件到管理器
}
