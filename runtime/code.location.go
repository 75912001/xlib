package runtime

import (
	xerror "github.com/75912001/xlib/error"
	xpool "github.com/75912001/xlib/pool"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
)

var ExecutablePath string // 程序所在路径(如为link,则为link所在的路径)
func init() {
	var err error
	if ExecutablePath, err = GetExecutablePath(); err != nil {
		panic(err)
	}
}

// 代码位置
type codeLocation struct {
	fileName string //文件名
	funcName string //函数名
	line     int    //行数
}

// String 信息
func (p *codeLocation) String() string {
	// 预计算字符串长度，避免 strings.Builder 动态扩容
	strLine := strconv.Itoa(p.line)
	size := len("file:") + len(p.fileName) +
		len(" line:") + len(strLine) +
		len(" func:") + len(p.funcName)

	b := xpool.Builder.Get()
	defer func() {
		xpool.Builder.Put(b)
	}()
	b.Grow(size) // 预分配内存空间

	b.WriteString("file:")
	b.WriteString(p.fileName)
	b.WriteString(" line:")
	b.WriteString(strLine)
	b.WriteString(" func:")
	b.WriteString(p.funcName)
	return b.String()
}

// Location 使用对象池复用 codeLocation 对象
var locationPool = sync.Pool{
	New: func() any {
		return &codeLocation{
			fileName: xerror.Unknown.Name(),
			funcName: xerror.Unknown.Name(),
		}
	},
}

// Location 获取代码位置
func Location() string {
	location := locationPool.Get().(*codeLocation)
	defer locationPool.Put(location)
	if pc, fileName, line, ok := runtime.Caller(1); ok {
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
