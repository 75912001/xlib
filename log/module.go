package log

import (
	"context"
)

type ILog interface {
	GetLevel() uint32
	SetLevel(level uint32) error
	Stop() error
	Trace(v ...interface{})
	TraceExtend(ctx context.Context, extendFields ExtendFields, v ...interface{})
	Tracef(format string, v ...interface{})
	TracefExtend(ctx context.Context, extendFields ExtendFields, format string, v ...interface{})
	Debug(v ...interface{})
	DebugLazy(vFunc func() []interface{})
	Debugf(format string, v ...interface{})
	DebugfLazy(formatFunc func() (string, []interface{}))
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}
