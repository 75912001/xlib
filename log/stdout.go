package log

import (
	xerror "github.com/75912001/xlib/error"
	"log"
	"os"
	"runtime"
	"time"
)

var stdOut = log.New(os.Stdout, "", 0)

// PrintInfo 输出到os.Stdout
func PrintInfo(v ...any) {
	funcName := xerror.Unknown.Name()
	pc, file, line, ok := runtime.Caller(calldepth1)
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	element := newEntry().
		withLevel(LevelInfo).
		withTime(time.Now()).
		withCallerInfo(line, file, funcName).
		withMessage("", v...)
	formatLogData(element)
	_ = stdOut.Output(calldepth2, element.outString)
}

// PrintfInfo 输出到os.Stdout
func PrintfInfo(format string, v ...any) {
	funcName := xerror.Unknown.Name()
	pc, file, line, ok := runtime.Caller(calldepth1)
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	element := newEntry().
		withLevel(LevelInfo).
		withTime(time.Now()).
		withCallerInfo(line, file, funcName).
		withMessage(format, v...)
	formatLogData(element)
	_ = stdOut.Output(calldepth2, element.outString)
}
