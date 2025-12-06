package control

type Event struct {
	ISwitch   ISwitchButton
	ICallBack ICallBack
}

// IOut 接口 导出
type IOut interface {
	Send(events ...any)
}
