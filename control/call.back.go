package control

type ICallBack interface {
	Execute() error             // 执行回调
	IParameters                 // 参数
	Clone(arg ...any) ICallBack // 克隆
}

type OnFunction func(args ...any) error
