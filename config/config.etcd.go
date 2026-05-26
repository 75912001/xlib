package config

import (
	"time"

	xerror "github.com/75912001/xlib/error"
	xetcdconstants "github.com/75912001/xlib/etcd/constants"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
)

type Etcd struct {
	Endpoints []string `yaml:"endpoints"` // etcd地址
	// YAML 须为 Go duration 字面量,如 100ms, 1s
	TTLDuration *time.Duration `yaml:"ttlDuration"` // ttlDuration [default]: xetcdconstants.TtlSecondDefault, e.g.:系统每10秒续约一次,该参数至少为11秒
}

func (p *Etcd) Configure() error {
	if len(p.Endpoints) == 0 {
		return errors.WithMessagef(xerror.Config, "endpoints is empty. %v", xruntime.Location())
	}
	if p.TTLDuration == nil {
		defaultValue := xetcdconstants.TtlDurationDefault
		p.TTLDuration = &defaultValue
	}
	return nil
}
