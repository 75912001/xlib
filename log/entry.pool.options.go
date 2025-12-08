package log

import (
	xcontrol "github.com/75912001/xlib/control"
	"sync"
)

// entry的内存池选项
type entryPoolOptions struct {
	poolSwitch   xcontrol.ISwitchButton // 内存池开关 [default]: true
	pool         *sync.Pool             // 内存池 [default]: &sync.Pool{New: func() any { return newEntry() }}
	newEntryFunc func() *entry          // 创建 entry 的方法 [default]: func() *entry { return p.pool.Get().(*entry) }
}

// newEntryPoolOptions 新的entryPoolOptions
func newEntryPoolOptions() *entryPoolOptions {
	pool := &sync.Pool{
		New: func() any {
			return newEntry()
		},
	}
	opt := &entryPoolOptions{
		poolSwitch: xcontrol.NewSwitchButton(true),
		pool:       pool,
		newEntryFunc: func() *entry {
			return pool.Get().(*entry)
		},
	}
	return opt
}

func (p *entryPoolOptions) merge(opts ...*entryPoolOptions) *entryPoolOptions {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.poolSwitch.IsOn() {
			p.poolSwitch.On()
		} else {
			p.poolSwitch.Off()
		}
		if opt.pool != nil {
			p.pool = opt.pool
		}
		if opt.newEntryFunc != nil {
			p.newEntryFunc = opt.newEntryFunc
		}
	}
	return p
}

// 配置
func (p *entryPoolOptions) configure() error {
	if p.poolSwitch.IsOn() {
		p.pool = &sync.Pool{
			New: func() any {
				return newEntry()
			},
		}
		p.newEntryFunc = func() *entry {
			return p.pool.Get().(*entry)
		}
	} else {
		p.newEntryFunc = func() *entry {
			return newEntry()
		}
	}
	return nil
}

// 将内存放回池中
func (p *entryPoolOptions) put(value *entry) {
	if p.poolSwitch.IsOn() {
		value.reset()
		p.pool.Put(value)
	}
}
