package runtime

import (
	xruntimeconstants "github.com/75912001/xlib/runtime/constants"
	"sync/atomic"
)

// 程序运行模式
var programRunMode atomic.Uint32

func init() {
	// 默认运行模式为发行模式
	programRunMode.Store(uint32(xruntimeconstants.RunModeRelease))
}
func SetRunMode(mode xruntimeconstants.RunMode) {
	programRunMode.Store(uint32(mode))
}

// IsDebug 是否为调试模式
func IsDebug() bool {
	return programRunMode.Load() == uint32(xruntimeconstants.RunModeDebug)
}

// IsRelease 是否为发行模式
func IsRelease() bool {
	return programRunMode.Load() == uint32(xruntimeconstants.RunModeRelease)
}
