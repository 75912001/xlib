package module

// IModule 模块
type IModule interface {
	Name() string
	Init(holder IHolder) error
	Start() error
	Stop() error
}

// IHolder 持有者
type IHolder interface {
	GetModule(t Type) IModule
}
