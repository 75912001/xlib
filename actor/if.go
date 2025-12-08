package actor

type IActor[Key comparable] interface {
	GetKey() Key
}

type IActorEvent interface {
	SendEvent(events ...*BehaviorEvent)
	SendEventWithResponse(event *BehaviorEvent) (response any, err error)
}
