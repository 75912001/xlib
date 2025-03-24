package exec

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type FuncAsyncStd func(data string) int

// CommandAsyncStdout 异步执行命令,并返回标准输出
// name:命令名称,默认为 /bin/bash
//
//	或者 "C:\\Program Files\\Git\\bin\\bash.exe"
//	args:"chmod +x /xx/xx/x.sh"
//
// args:"chmod +x /xx/xx/x.sh"
// funcStdout nil:disable stdout
// funcStderr nil:disable stderr
func CommandAsyncStdout(name string, args string, funcStdout FuncAsyncStd, funcStderr FuncAsyncStd) (err error) {
	if len(name) == 0 {
		name = "/bin/bash"
	}
	cmd := exec.Command(name, "-c", args)
	var stdout io.ReadCloser
	if stdout, err = cmd.StdoutPipe(); err != nil {
		return
	}
	var stderr io.ReadCloser
	if stderr, err = cmd.StderrPipe(); err != nil {
		return
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	if nil != funcStdout {
		go asyncStdout(stdout, funcStdout)
	}
	if nil != funcStderr {
		go asyncStdout(stderr, funcStderr)
	}

	if err := cmd.Wait(); err != nil {
		return err
	}
	return
}

func asyncStdout(reader io.ReadCloser, fun FuncAsyncStd) error {
	buf := make([]byte, 5)
	for {
		num, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if num > 0 {
			b := buf[:num]
			s := strings.Split(string(b), "\n")
			line := strings.Join(s[:len(s)-1], "\n")
			fun(fmt.Sprintln(line))
		}
	}
}
