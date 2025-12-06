package control

import (
	"unsafe"
)

type CallBack struct {
	onFunction OnFunction // 回调函数
	IParameters
}

func NewCallBack(onFunction OnFunction, arg ...any) *CallBack {
	newCallback := &CallBack{
		onFunction:  onFunction,
		IParameters: NewParameters(),
	}
	newCallback.IParameters.Override(arg...)
	return newCallback
}

// Clone 创建一个回调对象
func (p *CallBack) Clone(arg ...any) ICallBack {
	return NewCallBack(p.onFunction, arg...)
}

func (p *CallBack) Execute() error {
	if p.onFunction == nil {
		return nil
	}
	return p.onFunction(p.IParameters.Get()...)
}

// Equals 比较两个回调是否相等
//
//	[❗] 不适用闭包函数. 使用闭包函数比较时,创建新的函数对象,导致比较失败(编译器优化行为可能会掩盖该问题)
func (p *CallBack) Equals(other ICallBack) bool {
	if other == nil {
		return false
	}

	// 类型断言
	otherCallback, ok := other.(*CallBack)
	if !ok {
		return false
	}

	// 直接比较函数指针
	return *(*uintptr)(unsafe.Pointer(&p.onFunction)) == *(*uintptr)(unsafe.Pointer(&otherCallback.onFunction))
}
