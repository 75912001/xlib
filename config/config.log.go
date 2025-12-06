package config

import (
	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"path/filepath"
)

type Log struct {
	Level   *uint32 `yaml:"level"`   // 日志等级
	AbsPath *string `yaml:"absPath"` // 日志绝对路径		[default]: 当前执行的程序-绝对路径,指向启动当前进程的可执行文件-目录路径. e.g.:absPath/log
}

func (p *Log) Configure() error {
	if p.Level == nil {
		return errors.WithMessagef(xerror.Config, "level is nil. %v", xruntime.Location())
	}
	if p.AbsPath == nil {
		executablePath := filepath.Join(xruntime.ExecutablePath, "log")
		p.AbsPath = &executablePath
	}
	return nil
}
