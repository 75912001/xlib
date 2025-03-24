package control

// ISwitchButton 两种状态：开启、关闭
type ISwitchButton interface {
	On()         // 开
	Off()        // 关
	IsOn() bool  // 是否开
	IsOff() bool // 是否关
}
