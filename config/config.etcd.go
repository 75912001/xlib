package config

import (
	xerror "github.com/75912001/xlib/error"
	xetcdconstants "github.com/75912001/xlib/etcd/constants"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

type Etcd struct {
	Endpoints []string `yaml:"endpoints"` // etcd地址
	TTL       *int64   `yaml:"ttl"`       // ttl 秒		[default]: etcd.TtlSecondDefault 秒, e.g.:系统每10秒续约一次,该参数至少为11秒
}

func (p *Etcd) Configure() error {
	if len(p.Endpoints) == 0 {
		return errors.WithMessagef(xerror.Config, "endpoints is empty. %v", xruntime.Location())
	}
	if p.TTL == nil {
		defaultValue := xetcdconstants.TtlSecondDefault
		p.TTL = &defaultValue
	}
	return nil
}
