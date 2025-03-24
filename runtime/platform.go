package runtime

import "runtime"

// IsWindows win
func IsWindows() bool {
	return runtime.GOOS == `windows`
}

// IsLinux linux
func IsLinux() bool {
	return runtime.GOOS == `linux`
}

// IsDarwin darwin
func IsDarwin() bool {
	return `darwin` == runtime.GOOS
}
