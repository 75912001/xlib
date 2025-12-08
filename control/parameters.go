package control

// IParameters 参数
type IParameters interface {
	Override(args ...any) // 覆盖参数
	Get() []any           // 获取参数 [全部]
	Append(args ...any)   // 追加参数
}
