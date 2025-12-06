package config

import (
	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

type Nats struct {
	Name     *string  `yaml:"name"`     // 名称		[default]: "natsName"
	Addrs    []string `yaml:"addrs"`    // 地址
	User     *string  `yaml:"user"`     // 用户
	Password *string  `yaml:"password"` // 密码		[default]:"123456"
}

func (p *Nats) Configure() error {
	if p.Name == nil {
		defaultValue := "natsName"
		p.Name = &defaultValue
	}
	if len(p.Addrs) == 0 {
		return errors.WithMessagef(xerror.Config, "nats addrs is empty. %v", xruntime.Location())
	}
	if p.User == nil {
		return errors.WithMessagef(xerror.Config, "nats user is nil. %v", xruntime.Location())
	}
	if p.Password == nil {
		defaultValue := "123456"
		p.Password = &defaultValue
	}
	return nil
}
