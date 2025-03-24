package main

import (
	"context"
	"fmt"
	xlog "github.com/75912001/xlib/log"
	xruntime "github.com/75912001/xlib/runtime"
)

func logCallBackFunc(level uint32, outString string) {
	if xruntime.IsDebug() {
		fmt.Println("logCallBackFunc: ", level, outString)
	}
	return
}

func exampleLog() {
	if true {
		return
	}
	xruntime.SetRunMode(xruntime.RunModeDebug)
	fmt.Println("============================================================")
	xlog.PrintInfo("print info")
	xlog.PrintfInfo("print info %s", "format")
	xlog.PrintErr("print err")
	xlog.PrintfErr("print err %s", "format")
	fmt.Println("============================================================")
	var l xlog.ILog
	l, err := xlog.NewMgr(xlog.NewOptions().
		WithLevelCallBack(logCallBackFunc, xlog.LevelFatal, xlog.LevelError, xlog.LevelWarn),
	)
	if err != nil {
		panic(err)
	}
	xlog.PrintInfo("print info")
	xlog.PrintfInfo("print info %s", "format")
	xlog.PrintErr("print err")
	xlog.PrintfErr("print err %s", "format")

	l.Fatal("fatal")
	l.Fatalf("fatal %s", "format")
	l.Error("error")
	l.Errorf("error %s", "format")
	l.Warn("warn")
	l.Warnf("warn %s", "format")
	l.Info("info")
	l.Infof("info %s", "format")
	l.Debug("debug")
	{
		l.DebugLazy(func() []interface{} {
			return []interface{}{fmt.Sprintf("%v %v", "This is a complex log message", "msg")}
		})
	}
	l.Debugf("debug %s", "format")
	{
		l.DebugfLazy(func() (string, []interface{}) {
			return "format %v %v", []interface{}{"This is a complex log message", 111}
		})
	}
	l.Trace("trace")
	{
		ctx := context.Background()
		ctx = context.WithValue(ctx, xlog.TraceIDKey, "TraceIDKey.value1")
		ctx = context.WithValue(ctx, xlog.UserIDKey, uint64(668))
		l.TraceExtend(ctx, xlog.ExtendFields{"key1", "value1", 1001, 1}, "trace")
	}
	{
		ctx := context.Background()
		ctx = context.WithValue(ctx, xlog.TraceIDKey, "TraceIDKey.value1")
		l.TraceExtend(ctx, xlog.ExtendFields{"key1", "value1", 1001, 1, xlog.UserIDKey, uint64(7200)}, "trace")
	}
	l.Tracef("trace %s", "format")
	{
		ctx := context.Background()
		ctx = context.WithValue(ctx, xlog.TraceIDKey, "TraceIDKey.value1")
		ctx = context.WithValue(ctx, xlog.UserIDKey, uint64(668))
		l.TracefExtend(ctx, xlog.ExtendFields{"key1", "value1", 1001, 1}, "trace %s", "format")
	}
	err = l.Stop()
	if err != nil {
		panic(err)
	}

	return
}
