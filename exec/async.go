package exec

import (
	"fmt"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"io"
	"os/exec"
	"strings"
)

type FuncAsyncStd func(data string) int

// CommandAsyncStdout 异步执行命令,并返回标准输出
//
//	name:命令名称,默认为 /bin/bash
//	或者 "C:\\Program Files\\Git\\bin\\bash.exe"
//	args:"chmod +x /xx/xx/x.sh"
//	funcStdout nil:disable stdout
//	funcStderr nil:disable stderr
func CommandAsyncStdout(name string, args string, funcStdout FuncAsyncStd, funcStderr FuncAsyncStd) (err error) {
	if len(name) == 0 {
		name = "/bin/bash"
	}
	cmd := exec.Command(name, "-c", args)
	var stdout io.ReadCloser
	if stdout, err = cmd.StdoutPipe(); err != nil {
		return errors.WithMessagef(err, "get stdout pipe failed. %v", xruntime.Location())
	}
	var stderr io.ReadCloser
	if stderr, err = cmd.StderrPipe(); err != nil {
		return errors.WithMessagef(err, "get stderr pipe failed. %v", xruntime.Location())
	}
	if err = cmd.Start(); err != nil {
		return errors.WithMessagef(err, "start command failed. %v", xruntime.Location())
	}

	if nil != funcStdout {
		go func() {
			_ = asyncStdout(stdout, funcStdout)
		}()
	}
	if nil != funcStderr {
		go func() {
			_ = asyncStdout(stderr, funcStderr)
		}()
	}

	if err = cmd.Wait(); err != nil {
		return errors.WithMessagef(err, "wait command failed. %v", xruntime.Location())
	}
	return nil
}

func asyncStdout(reader io.ReadCloser, fun FuncAsyncStd) error {
	buf := make([]byte, 5)
	for {
		num, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return errors.WithMessagef(err, "read stdout failed. %v", xruntime.Location())
		}
		if num > 0 {
			b := buf[:num]
			s := strings.Split(string(b), "\n")
			line := strings.Join(s[:len(s)-1], "\n")
			fun(fmt.Sprintln(line))
		}
	}
}
