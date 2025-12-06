package log

import (
	xerror "github.com/75912001/xlib/error"
	"log"
	"os"
	"runtime"
	"time"
)

var stdErr = log.New(os.Stderr, "", 0)

// PrintErr 输出到os.Stderr
func PrintErr(v ...any) {
	funcName := xerror.Unknown.Name()
	pc, file, line, ok := runtime.Caller(calldepth1)
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	element := newEntry().
		withLevel(LevelError).
		withTime(time.Now()).
		withCallerInfo(line, file, funcName).
		withMessage("", v...)
	formatLogData(element)
	_ = stdErr.Output(calldepth2, element.outString)
}

// PrintfErr 输出到os.Stderr
func PrintfErr(format string, v ...any) {
	funcName := xerror.Unknown.Name()
	pc, file, line, ok := runtime.Caller(calldepth1)
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	element := newEntry().
		withLevel(LevelError).
		withTime(time.Now()).
		withCallerInfo(line, file, funcName).
		withMessage(format, v...)
	formatLogData(element)
	_ = stdErr.Output(calldepth2, element.outString)
}
