package config

import (
	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

type Redis struct {
	Name     *string  `yaml:"name"`     // 名称		[default]: "redisName"
	Addrs    []string `yaml:"addrs"`    // 地址
	Password *string  `yaml:"password"` // 密码		[default]:"123456"
}

func (p *Redis) Configure() error {
	if p.Name == nil {
		defaultValue := "redisName"
		p.Name = &defaultValue
	}
	if len(p.Addrs) == 0 {
		return errors.WithMessagef(xerror.Config, "redis addrs is empty. %v", xruntime.Location())
	}
	if p.Password == nil {
		defaultValue := "123456"
		p.Password = &defaultValue
	}
	return nil
}
