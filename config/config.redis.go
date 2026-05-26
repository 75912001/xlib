package config

import (
	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

type Redis struct {
	Name           *string  `yaml:"name"`           // 名称		[default]: "redisName"
	Addrs          []string `yaml:"addrs"`          // 地址
	Password       *string  `yaml:"password"`       // 密码		[default]:"123456"
	DialTimeoutMS  *int     `yaml:"dialTimeoutMS"`  // 连接超时时间-毫秒	[default]: 3000
	ReadTimeoutMS  *int     `yaml:"readTimeoutMS"`  // 读超时时间-毫秒	[default]: 3000
	WriteTimeoutMS *int     `yaml:"writeTimeoutMS"` // 写超时时间-毫秒	[default]: 3000
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
	if p.DialTimeoutMS == nil {
		defaultValue := 3000
		p.DialTimeoutMS = &defaultValue
	}
	if p.ReadTimeoutMS == nil {
		defaultValue := 3000
		p.ReadTimeoutMS = &defaultValue
	}
	if p.WriteTimeoutMS == nil {
		defaultValue := 3000
		p.WriteTimeoutMS = &defaultValue
	}
	return nil
}
