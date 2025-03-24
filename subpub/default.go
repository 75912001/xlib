package subpub

import (
	xcontrol "github.com/75912001/xlib/control"
	"github.com/pkg/errors"
)

type SubPub struct {
	subMap map[uint64][]xcontrol.ICallBack
}

func NewSubPub() *SubPub {
	return &SubPub{
		subMap: make(map[uint64][]xcontrol.ICallBack),
	}
}
func (p *SubPub) Subscribe(key uint64, onFunction func(...interface{}) error) error {
	p.subMap[key] = append(p.subMap[key], xcontrol.NewCallBack(onFunction))
	return nil
}

func (p *SubPub) Publish(key uint64, parameters ...interface{}) error {
	var err error
	if callbacks, exists := p.subMap[key]; exists {
		for _, callback := range callbacks {
			callback.Override(parameters...)
			if e := callback.Execute(); e != nil {
				if err == nil {
					err = e
				} else {
					err = errors.New(err.Error() + "; " + e.Error())
				}
			}
		}
	}
	return err
}
