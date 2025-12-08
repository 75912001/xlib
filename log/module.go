package log

import (
	"context"
)

var GLog ILog

type ILog interface {
	GetLevel() uint32
	SetLevel(level uint32) error
	Stop() error
	Trace(v ...any)
	TraceExtend(ctx context.Context, extendFields ExtendFields, v ...any)
	Tracef(format string, v ...any)
	TracefExtend(ctx context.Context, extendFields ExtendFields, format string, v ...any)
	Debug(v ...any)
	DebugExtend(ctx context.Context, extendFields ExtendFields, v ...any)
	DebugLazy(vFunc func() []any)
	Debugf(format string, v ...any)
	DebugfExtend(ctx context.Context, extendFields ExtendFields, format string, v ...any)
	DebugfLazy(formatFunc func() (string, []any))
	Info(v ...any)
	InfoExtend(ctx context.Context, extendFields ExtendFields, v ...any)
	Infof(format string, v ...any)
	InfofExtend(ctx context.Context, extendFields ExtendFields, format string, v ...any)
	Warn(v ...any)
	WarnExtend(ctx context.Context, extendFields ExtendFields, v ...any)
	Warnf(format string, v ...any)
	WarnfExtend(ctx context.Context, extendFields ExtendFields, format string, v ...any)
	Error(v ...any)
	ErrorExtend(ctx context.Context, extendFields ExtendFields, v ...any)
	Errorf(format string, v ...any)
	ErrorfExtend(ctx context.Context, extendFields ExtendFields, format string, v ...any)
	Fatal(v ...any)
	FatalExtend(ctx context.Context, extendFields ExtendFields, v ...any)
	Fatalf(format string, v ...any)
	FatalfExtend(ctx context.Context, extendFields ExtendFields, format string, v ...any)
}
