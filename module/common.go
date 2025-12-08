package module

type Type uint32

// CreatorFunction 定义模块创建器函数类型
//
//	将持有者传递给模块
type CreatorFunction func(holder IHolder) IModule
