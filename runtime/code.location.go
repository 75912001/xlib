package runtime

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	xerror "github.com/75912001/xlib/error"
	"github.com/pkg/errors"
)

var ExecutablePath string // 程序所在路径(如为link,则为link所在的路径)
func init() {
	var err error
	if ExecutablePath, err = GetExecutablePath(); err != nil {
		panic(err)
	}
}

// Location 获取代码位置
func Location() string {
	if pc, fileName, line, ok := runtime.Caller(1); ok {
		return "file:" + fileName + " line:" + strconv.Itoa(line) + " func:" + runtime.FuncForPC(pc).Name()
	}
	unknown := xerror.Unknown.Name()
	return "file:" + unknown + " line:0 func:" + unknown
}

// GetExecutablePath 获取当前执行的程序-绝对路径,指向启动当前进程的可执行文件-目录路径.
func GetExecutablePath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", errors.WithMessagef(err, "GetExecutablePath: %v", Location())
	}
	// 返回目录路径
	return filepath.Dir(path), nil
}

// GetExecutableName 获取当前执行的程序的名称
func GetExecutableName() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", errors.WithMessagef(err, "GetExecutableName: %v", Location())
	}
	// 返回程序名称
	return filepath.Base(path), nil
}

// GetRealExecutablePath 获取当前执行的程序,符号链接的实际路径-绝对路径,指向启动当前进程的可执行文件,符号链接的实际路径-目录路径.
func GetRealExecutablePath() (currentPath string, err error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", errors.WithMessagef(err, "GetRealExecutablePath: %v", Location())
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return "", errors.WithMessagef(err, "GetRealExecutablePath: %v", Location())
	}
	return filepath.Dir(exePath), nil
}
