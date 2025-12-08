package module

type Module struct {
	Type            Type            // 模块类型
	Name            string          // 模块名称
	IModule         IModule         // 模块
	CreatorFunction CreatorFunction // 模块创建器
}
