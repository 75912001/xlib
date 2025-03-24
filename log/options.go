package log

import (
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

// options contains options to configure a server mgrInstance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type options struct {
	level            *uint32           // 日志等级允许的最小等级 [default]: LevelOn
	absPath          *string           // 日志绝对路径 [default]: 当前执行的程序-绝对路径,指向启动当前进程的可执行文件-目录路径. e.g.:absPath/log
	isReportCaller   *bool             // 是否打印调用信息 [default]: true
	namePrefix       *string           // 日志名 前缀 [default]: 当前执行的程序名称
	isWriteFile      *bool             // 是否写文件 [default]: true
	entryPoolOptions *entryPoolOptions // entry的内存池选项 [default]: newEntryPoolOptions()
	levelSubscribe   *levelSubscribe   // 日志结果的回调. 该执行过程在 log 的写过程中, 位于写 log 的 goroutine 中 [default]: nil
}

// NewOptions 新的Options
func NewOptions() *options {
	return &options{}
}

func (p *options) WithLevel(level uint32) *options {
	p.level = &level
	return p
}

func (p *options) WithAbsPath(absPath string) *options {
	p.absPath = &absPath
	return p
}

func (p *options) WithIsReportCaller(isReportCaller bool) *options {
	p.isReportCaller = &isReportCaller
	return p
}

func (p *options) WithNamePrefix(namePrefix string) *options {
	p.namePrefix = &namePrefix
	return p
}

func (p *options) WithIsWriteFile(isWriteFile bool) *options {
	p.isWriteFile = &isWriteFile
	return p
}

func (p *options) WithEntryPoolOptions(entryPoolOptions *entryPoolOptions) *options {
	p.entryPoolOptions = entryPoolOptions
	return p
}

func (p *options) WithLevelCallBack(callbackFunc CallBackFunc, subLevel ...uint32) *options {
	p.levelSubscribe = newLevelSubscribe()
	p.levelSubscribe.callBackFunc = callbackFunc
	for _, level := range subLevel {
		p.levelSubscribe.subMap[level] = struct{}{}
	}
	return p
}

// merge combines the given *options into a single *options in a last one wins fashion.
// The specified options are merged with the existing options on the server, with the specified options taking
// precedence.
func (p *options) merge(opts ...*options) *options {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.level != nil {
			p.level = opt.level
		}
		if opt.absPath != nil {
			p.absPath = opt.absPath
		}
		if opt.isReportCaller != nil {
			p.isReportCaller = opt.isReportCaller
		}
		if opt.namePrefix != nil {
			p.namePrefix = opt.namePrefix
		}
		if opt.isWriteFile != nil {
			p.isWriteFile = opt.isWriteFile
		}
		if opt.entryPoolOptions != nil {
			p.entryPoolOptions = p.entryPoolOptions.merge(opt.entryPoolOptions)
		}
		if opt.levelSubscribe != nil {
			p.levelSubscribe = opt.levelSubscribe
		}
	}
	return p
}

// 配置
func (p *options) configure() error {
	if p.level == nil {
		var level = LevelOn
		p.level = &level
	}
	if p.absPath == nil {
		executablePath, err := xruntime.GetExecutablePath()
		if err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
		executablePath = filepath.Join(executablePath, "log")
		p.absPath = &executablePath
	}
	if err := os.MkdirAll(*p.absPath, os.ModePerm); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	if p.isReportCaller == nil {
		var reportCaller = true
		p.isReportCaller = &reportCaller
	}
	if p.namePrefix == nil {
		executableName, err := xruntime.GetExecutableName()
		if err != nil {
			return errors.WithMessage(err, xruntime.Location())
		}
		p.namePrefix = &executableName
	}
	if p.isWriteFile == nil {
		var writeFile = true
		p.isWriteFile = &writeFile
	}
	if p.entryPoolOptions == nil {
		p.entryPoolOptions = newEntryPoolOptions()
	}
	if err := p.entryPoolOptions.configure(); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}
