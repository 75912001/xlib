package config

import (
	"time"

	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

type Redis struct {
	Name     *string  `yaml:"name"`     // 名称		[default]: "redisName"
	Addrs    []string `yaml:"addrs"`    // 地址
	Password *string  `yaml:"password"` // 密码		[default]:"123456"
	// YAML 须为 Go duration 字面量,如 100ms, 1s
	DialTimeout *time.Duration `yaml:"dialTimeout"` // 连接超时时间	[default]: 3s
	// YAML 须为 Go duration 字面量,如 100ms, 1s
	ReadTimeout *time.Duration `yaml:"readTimeout"` // 读超时时间	[default]: 3s
	// YAML 须为 Go duration 字面量,如 100ms, 1s
	WriteTimeout *time.Duration `yaml:"writeTimeout"` // 写超时时间	[default]: 3s
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
	if p.DialTimeout == nil {
		defaultValue := time.Second * 3
		p.DialTimeout = &defaultValue
	}
	if p.ReadTimeout == nil {
		defaultValue := time.Second * 3
		p.ReadTimeout = &defaultValue
	}
	if p.WriteTimeout == nil {
		defaultValue := time.Second * 3
		p.WriteTimeout = &defaultValue
	}
	return nil
}
