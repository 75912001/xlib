package unuse

import (
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Sub struct {
	*Jaeger  `yaml:"jaeger"`  // todo
	*MongoDB `yaml:"mongoDB"` // todo
}

func (p *Sub) Unmarshal(strYaml string) error {
	if err := yaml.Unmarshal([]byte(strYaml), &p); err != nil {
		return errors.WithMessage(err, xruntime.Location())
	}
	return nil
}

type Jaeger struct {
	Addrs []string `yaml:"addrs"`
}
