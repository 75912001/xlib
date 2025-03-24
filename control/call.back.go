package control

type ICallBack interface {
	Execute() error // 执行回调
	IParameters     // 参数
}
