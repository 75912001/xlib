package control

// IParameters 参数
type IParameters interface {
	Override(args ...any) // 覆盖参数
	// Get 获取参数 [全部]
	// 约定(仅注释约束，调用方须遵守): 禁止修改返回值,禁止长期持有返回值并在 Override,Append 等之后仍当作最新参数使用,须重新调用 Get
	Get() []any
	Append(args ...any) // 追加参数
}
