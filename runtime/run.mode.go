package runtime

import (
	xruntimeconstants "github.com/75912001/xlib/runtime/constants"
)

// 程序运行模式
var programRunMode = xruntimeconstants.RunModeRelease

func SetRunMode(mode xruntimeconstants.RunMode) {
	programRunMode = mode
}

// IsDebug 是否为调试模式
func IsDebug() bool {
	return programRunMode == xruntimeconstants.RunModeDebug
}

// IsRelease 是否为发行模式
func IsRelease() bool {
	return programRunMode == xruntimeconstants.RunModeRelease
}
