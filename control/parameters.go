package control

// IParameters 参数
type IParameters interface {
	Override(parameters ...interface{}) // 覆盖参数
	Get() []interface{}                 // 获取参数 [全部]
	Append(parameters ...interface{})   // 追加参数
}
