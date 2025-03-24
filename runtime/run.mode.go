package runtime

// runMode 运行模式
type runMode uint32

const (
	RunModeRelease runMode = 0 // release 模式
	RunModeDebug   runMode = 1 // debug 模式
)

// 程序运行模式
var programRunMode = RunModeRelease

func SetRunMode(mode runMode) {
	programRunMode = mode
}

// IsDebug 是否为调试模式
func IsDebug() bool {
	return programRunMode == RunModeDebug
}

// IsRelease 是否为发行模式
func IsRelease() bool {
	return programRunMode == RunModeRelease
}
