package runtime

import (
	"fmt"
	xerror "github.com/75912001/xlib/error"
	"os"
	"path/filepath"
	"runtime"
)

// 代码位置
type codeLocation struct {
	fileName string //文件名
	funcName string //函数名
	line     int    //行数
}

// String 信息
func (p *codeLocation) String() string {
	return fmt.Sprintf("file:%v line:%v func:%v", p.fileName, p.line, p.funcName)
}

// Location 获取代码位置
func Location() string {
	location := &codeLocation{
		fileName: xerror.Unknown.Name(),
		funcName: xerror.Unknown.Name(),
	}
	pc, fileName, line, ok := runtime.Caller(1)
	if ok {
		location.fileName = fileName
		location.line = line
		location.funcName = runtime.FuncForPC(pc).Name()
	}
	return location.String()
}

// GetExecutablePath 获取当前执行的程序-绝对路径,指向启动当前进程的可执行文件-目录路径.
func GetExecutablePath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	// 返回目录路径
	return filepath.Dir(path), nil
}

// GetExecutableName 获取当前执行的程序的名称
func GetExecutableName() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	// 返回程序名称
	return filepath.Base(path), nil
}

// GetRealExecutablePath 获取当前执行的程序,符号链接的实际路径-绝对路径,指向启动当前进程的可执行文件,符号链接的实际路径-目录路径.
func GetRealExecutablePath() (currentPath string, err error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}
