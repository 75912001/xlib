// 日志
// 使用系统log,自带锁
// 使用协程操作io输出日志
// release 每小时自动创建新的日志文件
// debug 每天自动创建新的日志文件

package log

import (
	"context"
	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	xtime "github.com/75912001/xlib/time"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

// NewMgr 创建日志管理器
func NewMgr(opts ...*options) (ILog, error) {
	m := &mgr{}
	err := m.handleOptions(opts...)
	if err != nil {
		return nil, errors.WithMessagef(err, "handle options failed. %v", xruntime.Location())
	}
	err = m.start()
	if err != nil {
		return nil, errors.WithMessagef(err, "start failed. %v", xruntime.Location())
	}
	return m, nil
}

// 日志管理器
type mgr struct {
	options         *options
	loggerSlice     [LevelOn]*log.Logger // 日志实例 [note]:使用时,注意协程安全
	logChan         chan *entry          // 日志写入通道
	waitGroupOutPut sync.WaitGroup       // 同步锁 用于日志退出时,等待完全输出
	logDuration     int                  // 日志分割刻度,变化时,使用新的日志文件 按天或者小时  e.g.: 20210819 或 2021081901
	openFiles       []*os.File           // 当前打开的文件
	timeMgr         *xtime.Mgr
}

func (p *mgr) handleOptions(opts ...*options) error {
	p.options = NewOptions().merge(opts...)
	if err := p.options.configure(); err != nil {
		return errors.WithMessagef(err, "configure failed. %v", xruntime.Location())
	}
	return nil
}

// 开始
func (p *mgr) start() error {
	// 初始化logger
	for i := LevelOff; i < LevelOn; i++ {
		p.loggerSlice[i] = log.New(os.Stdout, "", 0)
	}
	p.logChan = make(chan *entry, logChannelEntryCapacity)
	p.timeMgr = xtime.NewMgr()
	// 初始化各级别的日志输出
	if err := newWriters(p); err != nil {
		return errors.WithMessagef(err, "new writers failed. %v", xruntime.Location())
	}
	p.waitGroupOutPut.Add(1)
	go func() {
		defer func() {
			if xruntime.IsRelease() {
				if err := recover(); err != nil {
					PrintErr(p, xerror.GoroutinePanic, err, string(debug.Stack()))
				}
			}
			p.waitGroupOutPut.Done()
			PrintInfo(xerror.GoroutineDone)
		}()
		doLog(p)
	}()
	return nil
}

// GetLevel 获取日志等级
func (p *mgr) GetLevel() uint32 {
	return *p.options.level
}

// getLogDuration 取得日志刻度
func (p *mgr) getLogDuration(sec int64) int {
	var logFormat string
	if xruntime.IsRelease() {
		logFormat = "2006010215" //年月日小时
	} else {
		logFormat = "20060102" //年月日
	}
	durationStr := time.Unix(sec, 0).Format(logFormat)
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		PrintfErr("strconv.Atoi sec:%v durationStr:%v err:%v", sec, durationStr, err)
	}
	return duration
}

// doLog 处理日志
func doLog(p *mgr) {
	for v := range p.logChan {
		formatLogData(v)
		p.callBack(v)
		// 检查自动切换日志
		if p.logDuration != p.getLogDuration(v.time.Unix()) {
			if err := newWriters(p); err != nil {
				PrintfErr("log duration changed, init writers failed, err:%v", err)
				p.options.entryPoolOptions.put(v)
				continue
			}
		}
		if *p.options.isWriteFile {
			p.loggerSlice[v.level].Print(v.outString)
		}
		p.options.entryPoolOptions.put(v)
	}
	// goroutine 退出,再设置chan为nil, (如果没有退出就设置为nil, 读chan == nil  会 block)
	p.logChan = nil
}

// SetLevel 设置日志等级
func (p *mgr) SetLevel(level uint32) error {
	if level < LevelOff || LevelOn < level {
		return errors.WithMessagef(xerror.Level, "level is invalid. %v", xruntime.Location())
	}
	p.options.WithLevel(level)
	return nil
}

// newWriters 初始化各级别的日志输出
func newWriters(p *mgr) error {
	// 检查是否要关闭文件
	for i := range p.openFiles {
		if err := p.openFiles[i].Close(); err != nil {
			return errors.WithMessagef(err, "close file failed. %v", xruntime.Location())
		}
	}
	second := p.timeMgr.NowTime().Unix()
	logDuration := p.getLogDuration(second)
	normalWriter, err := newNormalFileWriter(*p.options.absPath, *p.options.namePrefix, logDuration)
	if err != nil {
		return errors.WithMessagef(err, "new normal file writer failed. %v", xruntime.Location())
	}
	errorWriter, err := newErrorFileWriter(*p.options.absPath, *p.options.namePrefix, logDuration)
	if err != nil {
		return errors.WithMessagef(err, "new error file writer failed. %v", xruntime.Location())
	}
	p.logDuration = logDuration
	allWriter := io.MultiWriter(normalWriter, errorWriter)
	// 标准输出,标准错误重定向
	stdOut.SetOutput(normalWriter)
	stdErr.SetOutput(allWriter)

	p.loggerSlice[LevelFatal].SetOutput(allWriter)
	p.loggerSlice[LevelError].SetOutput(allWriter)
	p.loggerSlice[LevelWarn].SetOutput(allWriter)
	p.loggerSlice[LevelInfo].SetOutput(normalWriter)
	p.loggerSlice[LevelDebug].SetOutput(normalWriter)
	p.loggerSlice[LevelTrace].SetOutput(normalWriter)
	// 记录打开的文件
	p.openFiles = p.openFiles[0:0]
	p.openFiles = append(p.openFiles, normalWriter)
	p.openFiles = append(p.openFiles, errorWriter)
	return nil
}

// Stop 停止
func (p *mgr) Stop() error {
	if p.logChan != nil {
		// close chan, for range 读完chan会退出.
		close(p.logChan)
		// 等待logChan 的for range 退出.
		p.waitGroupOutPut.Wait()
	}
	// 检查是否要关闭文件
	if len(p.openFiles) > 0 {
		for i := range p.openFiles {
			_ = p.openFiles[i].Close()
		}
		p.openFiles = p.openFiles[0:0]
	}
	return nil
}

// callBack 处理回调
func (p *mgr) callBack(entry *entry) {
	if p.options.levelSubscribe == nil {
		return
	}
	if !p.options.levelSubscribe.isSubscribe(entry.level) {
		return
	}
	p.options.levelSubscribe.callBack.Override(entry.level, entry.outString)
	_ = p.options.levelSubscribe.callBack.Execute()
}

func (p *mgr) newEntry() *entry {
	return p.options.entryPoolOptions.newEntryFunc()
}

// log 记录日志
func (p *mgr) log(entry *entry, level uint32, v ...any) {
	entry.withLevel(level).
		withTime(p.timeMgr.NowTime()).
		withMessage("", v...)
	if *p.options.isReportCaller {
		pc, file, line, ok := runtime.Caller(calldepth2)
		funcName := xerror.Unknown.Name()
		if !ok {
			line = 0
		} else {
			funcName = runtime.FuncForPC(pc).Name()
		}
		entry.withCallerInfo(line, file, funcName)
	}
	p.logChan <- entry
}

// logf 记录日志
func (p *mgr) logf(entry *entry, level uint32, format string, v ...any) {
	entry.withLevel(level).
		withTime(p.timeMgr.NowTime()).
		withMessage(format, v...)
	if *p.options.isReportCaller {
		pc, file, line, ok := runtime.Caller(calldepth2)
		funcName := xerror.Unknown.Name()
		if !ok {
			line = 0
		} else {
			funcName = runtime.FuncForPC(pc).Name()
		}
		entry.withCallerInfo(line, file, funcName)
	}
	p.logChan <- entry
}

// Trace 踪迹日志
func (p *mgr) Trace(v ...any) {
	if p.GetLevel() < LevelTrace {
		return
	}
	p.log(p.newEntry(), LevelTrace, v...)
}

func (p *mgr) TraceExtend(ctx context.Context, extendFields ExtendFields, v ...any) {
	if p.GetLevel() < LevelTrace {
		return
	}
	element := p.newEntry()
	element.WithContext(ctx).WithExtendFields(extendFields)
	p.log(element, LevelTrace, v...)
}

// Tracef 踪迹日志
func (p *mgr) Tracef(format string, v ...any) {
	if p.GetLevel() < LevelTrace {
		return
	}
	p.logf(p.newEntry(), LevelTrace, format, v...)
}

func (p *mgr) TracefExtend(ctx context.Context, extendFields ExtendFields, format string, v ...any) {
	if p.GetLevel() < LevelTrace {
		return
	}
	element := p.newEntry()
	element.WithContext(ctx).WithExtendFields(extendFields)
	p.logf(element, LevelTrace, format, v...)
}

// Debug 调试日志
func (p *mgr) Debug(v ...any) {
	if p.GetLevel() < LevelDebug {
		return
	}
	p.log(p.newEntry(), LevelDebug, v...)
}

func (p *mgr) DebugExtend(ctx context.Context, extendFields ExtendFields, v ...any) {
	if p.GetLevel() < LevelDebug {
		return
	}
	element := p.newEntry()
	element.WithContext(ctx).WithExtendFields(extendFields)
	p.log(element, LevelDebug, v...)
}

// DebugLazy 调试日志-惰性
//
//	等级满足之后才会计算
func (p *mgr) DebugLazy(vFunc func() []any) {
	if p.GetLevel() < LevelDebug {
		return
	}
	v := vFunc()
	p.log(p.newEntry(), LevelDebug, v...)
}

// Debugf 调试日志
func (p *mgr) Debugf(format string, v ...any) {
	if p.GetLevel() < LevelDebug {
		return
	}
	p.logf(p.newEntry(), LevelDebug, format, v...)
}

func (p *mgr) DebugfExtend(ctx context.Context, extendFields ExtendFields, format string, v ...any) {
	if p.GetLevel() < LevelDebug {
		return
	}
	element := p.newEntry()
	element.WithContext(ctx).WithExtendFields(extendFields)
	p.logf(element, LevelDebug, format, v...)
}

// DebugfLazy 调试日志-惰性
//
//	等级满足之后才会计算
func (p *mgr) DebugfLazy(formatFunc func() (string, []any)) {
	if p.GetLevel() < LevelDebug {
		return
	}
	format, v := formatFunc()
	p.logf(p.newEntry(), LevelDebug, format, v...)
}

// Info 信息日志
func (p *mgr) Info(v ...any) {
	if p.GetLevel() < LevelInfo {
		return
	}
	p.log(p.newEntry(), LevelInfo, v...)
}

func (p *mgr) InfoExtend(ctx context.Context, extendFields ExtendFields, v ...any) {
	if p.GetLevel() < LevelInfo {
		return
	}
	element := p.newEntry()
	element.WithContext(ctx).WithExtendFields(extendFields)
	p.log(element, LevelInfo, v...)
}

// Infof 信息日志
func (p *mgr) Infof(format string, v ...any) {
	if p.GetLevel() < LevelInfo {
		return
	}
	p.logf(p.newEntry(), LevelInfo, format, v...)
}

func (p *mgr) InfofExtend(ctx context.Context, extendFields ExtendFields, format string, v ...any) {
	if p.GetLevel() < LevelInfo {
		return
	}
	element := p.newEntry()
	element.WithContext(ctx).WithExtendFields(extendFields)
	p.logf(element, LevelInfo, format, v...)
}

// Warn 警告日志
func (p *mgr) Warn(v ...any) {
	if p.GetLevel() < LevelWarn {
		return
	}
	p.log(p.newEntry(), LevelWarn, v...)
}

func (p *mgr) WarnExtend(ctx context.Context, extendFields ExtendFields, v ...any) {
	if p.GetLevel() < LevelWarn {
		return
	}
	element := p.newEntry()
	element.WithContext(ctx).WithExtendFields(extendFields)
	p.log(element, LevelWarn, v...)
}

// Warnf 警告日志
func (p *mgr) Warnf(format string, v ...any) {
	if p.GetLevel() < LevelWarn {
		return
	}
	p.logf(p.newEntry(), LevelWarn, format, v...)
}

func (p *mgr) WarnfExtend(ctx context.Context, extendFields ExtendFields, format string, v ...any) {
	if p.GetLevel() < LevelWarn {
		return
	}
	element := p.newEntry()
	element.WithContext(ctx).WithExtendFields(extendFields)
	p.logf(element, LevelWarn, format, v...)
}

// Error 错误日志
func (p *mgr) Error(v ...any) {
	if p.GetLevel() < LevelError {
		return
	}
	p.log(p.newEntry(), LevelError, v...)
}

func (p *mgr) ErrorExtend(ctx context.Context, extendFields ExtendFields, v ...any) {
	if p.GetLevel() < LevelError {
		return
	}
	element := p.newEntry()
	element.WithContext(ctx).WithExtendFields(extendFields)
	p.log(element, LevelError, v...)
}

// Errorf 错误日志
func (p *mgr) Errorf(format string, v ...any) {
	if p.GetLevel() < LevelError {
		return
	}
	p.logf(p.newEntry(), LevelError, format, v...)
}

func (p *mgr) ErrorfExtend(ctx context.Context, extendFields ExtendFields, format string, v ...any) {
	if p.GetLevel() < LevelError {
		return
	}
	element := p.newEntry()
	element.WithContext(ctx).WithExtendFields(extendFields)
	p.logf(element, LevelError, format, v...)
}

// Fatal 致命日志
func (p *mgr) Fatal(v ...any) {
	if p.GetLevel() < LevelFatal {
		return
	}
	p.log(p.newEntry(), LevelFatal, v...)
}

func (p *mgr) FatalExtend(ctx context.Context, extendFields ExtendFields, v ...any) {
	if p.GetLevel() < LevelFatal {
		return
	}
	element := p.newEntry()
	element.WithContext(ctx).WithExtendFields(extendFields)
	p.log(element, LevelFatal, v...)
}

// Fatalf 致命日志
func (p *mgr) Fatalf(format string, v ...any) {
	if p.GetLevel() < LevelFatal {
		return
	}
	p.logf(p.newEntry(), LevelFatal, format, v...)
}

func (p *mgr) FatalfExtend(ctx context.Context, extendFields ExtendFields, format string, v ...any) {
	if p.GetLevel() < LevelFatal {
		return
	}
	element := p.newEntry()
	element.WithContext(ctx).WithExtendFields(extendFields)
	p.logf(element, LevelFatal, format, v...)
}
