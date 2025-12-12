package actor

type IActor[Key comparable] interface {
	GetKey() Key
}

type IActorMsg interface {
	SendMsg(msg ...*Msg)
	SendMsgAsync(msg *Msg) (resp any, err error)
}
