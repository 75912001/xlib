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
	DialTimeoutDuration *time.Duration `yaml:"dialTimeoutDuration"` // 连接超时时间	[default]: 3s
	// YAML 须为 Go duration 字面量,如 100ms, 1s
	ReadTimeoutDuration *time.Duration `yaml:"readTimeoutDuration"` // 读超时时间	[default]: 3s
	// YAML 须为 Go duration 字面量,如 100ms, 1s
	WriteTimeoutDuration *time.Duration `yaml:"writeTimeoutDuration"` // 写超时时间	[default]: 3s
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
	if p.DialTimeoutDuration == nil {
		defaultValue := time.Second * 3
		p.DialTimeoutDuration = &defaultValue
	}
	if p.ReadTimeoutDuration == nil {
		defaultValue := time.Second * 3
		p.ReadTimeoutDuration = &defaultValue
	}
	if p.WriteTimeoutDuration == nil {
		defaultValue := time.Second * 3
		p.WriteTimeoutDuration = &defaultValue
	}
	return nil
}
