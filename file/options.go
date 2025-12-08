package file

// options 选项
type options struct {
	overwrite bool // 覆盖
	append    bool // 追加
}

// NewOptions 创建选项
func NewOptions() *options {
	return &options{
		overwrite: true,  // 默认覆盖
		append:    false, // 默认不追加
	}
}

// Overwrite 覆盖
func (fo *options) Overwrite() *options {
	fo.overwrite = true
	fo.append = false
	return fo
}

// Append 追加
func (fo *options) Append() *options {
	fo.append = true
	fo.overwrite = false
	return fo
}
