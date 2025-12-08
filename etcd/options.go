package etcd

import (
	xcontrol "github.com/75912001/xlib/control"
	xerror "github.com/75912001/xlib/error"
	xetcdconstants "github.com/75912001/xlib/etcd/constants"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"time"
)

type Options struct {
	endpoints            []string           // 地址
	ttl                  *int64             // Time To Live, etcd内部会按照 ttl/3 的时间(最小1秒),保持连接
	grantLeaseMaxRetries *int               // 授权租约 最大 重试次数 [default: grantLeaseMaxRetriesDefault]
	dialTimeout          *time.Duration     // dialTimeout is the timeout for failing to establish a connection. [default: dialTimeoutDefault]
	iOut                 xcontrol.IOut      // 传出 [nil: 则不传出事件] // 当设置了 watchKeyPrefix 时, 该接口用于传出事件,与 watchKeyPrefix 同时生效
	watchKeyPrefix       *string            // 监视的键前缀  [nil: 则不监视] // 不可与 key 同时设置为 nil
	key                  *string            // 本服务的 etcd key [nil: 则不设置key] // 不可与 watchKeyPrefix 同时设置为 nil
	AddCallback          xcontrol.ICallBack // 增加回调 [default: nil]
	UpdateCallback       xcontrol.ICallBack // 更新回调 [default: nil]
	DelCallback          xcontrol.ICallBack // 删除回调 [default: nil]
}

// NewOptions 新的Options
func NewOptions() *Options {
	return &Options{}
}

func (p *Options) WithAddCallback(callback xcontrol.ICallBack) *Options {
	p.AddCallback = callback
	return p
}

func (p *Options) WithUpdateCallback(callback xcontrol.ICallBack) *Options {
	p.UpdateCallback = callback
	return p
}

func (p *Options) WithDelCallback(callback xcontrol.ICallBack) *Options {
	p.DelCallback = callback
	return p
}

func (p *Options) WithEndpoints(endpoints []string) *Options {
	p.endpoints = p.endpoints[0:0]
	p.endpoints = append(p.endpoints, endpoints...)
	return p
}

func (p *Options) WithTTL(ttl int64) *Options {
	p.ttl = &ttl
	return p
}

func (p *Options) WithGrantLeaseMaxRetries(retries int) *Options {
	p.grantLeaseMaxRetries = &retries
	return p
}

func (p *Options) WithDialTimeout(dialTimeout time.Duration) *Options {
	p.dialTimeout = &dialTimeout
	return p
}

func (p *Options) WithIOut(iOut xcontrol.IOut) *Options {
	p.iOut = iOut
	return p
}

func (p *Options) WithWatchKeyPrefix(watchKeyPrefix string) *Options {
	p.watchKeyPrefix = &watchKeyPrefix
	return p
}

func (p *Options) GetKey() string {
	if p.key == nil {
		return ""
	}
	return *p.key
}

func (p *Options) WithKey(key string) *Options {
	p.key = &key
	return p
}

func MergeOptions(opts ...*Options) *Options {
	no := NewOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if len(opt.endpoints) != 0 {
			no.WithEndpoints(opt.endpoints)
		}
		if opt.ttl != nil {
			no.WithTTL(*opt.ttl)
		}
		if opt.grantLeaseMaxRetries != nil {
			no.WithGrantLeaseMaxRetries(*opt.grantLeaseMaxRetries)
		}
		if opt.dialTimeout != nil {
			no.WithDialTimeout(*opt.dialTimeout)
		}
		if opt.iOut != nil {
			no.WithIOut(opt.iOut)
		}
		if opt.watchKeyPrefix != nil {
			no.WithWatchKeyPrefix(*opt.watchKeyPrefix)
		}
		if opt.key != nil {
			no.WithKey(*opt.key)
		}
		if opt.AddCallback != nil {
			no.WithAddCallback(opt.AddCallback)
		}
		if opt.UpdateCallback != nil {
			no.WithUpdateCallback(opt.UpdateCallback)
		}
		if opt.DelCallback != nil {
			no.WithDelCallback(opt.DelCallback)
		}
	}
	return no
}

// 配置
func configure(opts *Options) error {
	if len(opts.endpoints) == 0 {
		return errors.WithMessagef(xerror.Param, "endpoints is empty. %v", xruntime.Location())
	}
	if opts.ttl == nil {
		return errors.WithMessagef(xerror.Param, "ttl is nil. %v", xruntime.Location())
	}
	if opts.grantLeaseMaxRetries == nil {
		var v = xetcdconstants.GrantLeaseMaxRetriesDefault
		opts.grantLeaseMaxRetries = &v
	}
	if opts.dialTimeout == nil {
		opts.WithDialTimeout(xetcdconstants.DialTimeoutDefault)
	}
	if opts.key == nil && opts.watchKeyPrefix == nil { // 既不监视 也 不设置key
		return errors.WithMessagef(xerror.Param, "key and watchKeyPrefix are nil. %v", xruntime.Location())
	}
	if (opts.watchKeyPrefix != nil && opts.iOut == nil) ||
		(opts.iOut != nil && opts.watchKeyPrefix == nil) { // 要么同时为 nil, 要么同时不为 nil
		return errors.WithMessagef(xerror.Param, "watchKeyPrefix,iOut 要么同时为 nil, 要么同时不为 nil. %v", xruntime.Location())
	}
	return nil
}
