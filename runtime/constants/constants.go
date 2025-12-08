package constants

// RunMode 运行模式
type RunMode uint32

const (
	RunModeRelease RunMode = 0 // release 模式
	RunModeDebug   RunMode = 1 // debug 模式
)
