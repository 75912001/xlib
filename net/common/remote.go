package common

type IRemote interface {
	ISend
	IsConnect() bool
	Start(connOptions *ConnOptions, event IEvent, handler IHandler)
	Stop()
	GetIP() string
	GetDisconnectReason() DisconnectReason
	SetDisconnectReason(reason DisconnectReason)
}
