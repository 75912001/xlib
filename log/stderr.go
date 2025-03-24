package log

import (
	"fmt"
	xerror "github.com/75912001/xlib/error"
	"log"
	"os"
	"runtime"
	"time"
)

var stdErr = log.New(os.Stderr, "", 0)

// PrintErr 输出到os.Stderr
func PrintErr(v ...interface{}) {
	funcName := xerror.Unknown.Name()
	pc, _, line, ok := runtime.Caller(calldepth1)
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	element := newEntry().
		withLevel(LevelError).
		withTime(time.Now()).
		withCallerInfo(fmt.Sprintf(callerInfoFormat, line, funcName)).
		withMessage(fmt.Sprint(v...))
	formatLogData(element)
	_ = stdErr.Output(calldepth2, element.outString)
}

// PrintfErr 输出到os.Stderr
func PrintfErr(format string, v ...interface{}) {
	funcName := xerror.Unknown.Name()
	pc, _, line, ok := runtime.Caller(calldepth1)
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}
	element := newEntry().
		withLevel(LevelError).
		withTime(time.Now()).
		withCallerInfo(fmt.Sprintf(callerInfoFormat, line, funcName)).
		withMessage(fmt.Sprintf(format, v...))
	formatLogData(element)
	_ = stdErr.Output(calldepth2, element.outString)
}
