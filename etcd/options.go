package etcd

import (
	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"time"
)

type options struct {
	addrs                []string           // 地址
	ttl                  *int64             // Time To Live, etcd内部会按照 ttl/3 的时间(最小1秒),保持连接
	grantLeaseMaxRetries *int               // 授权租约 最大 重试次数 [default: grantLeaseMaxRetriesDefault]
	dialTimeout          *time.Duration     // dialTimeout is the timeout for failing to establish a connection. [default: dialTimeoutDefault]
	eventChan            chan<- interface{} // 传出 channel
	watchKeyPrefix       *string            // 监视的键前缀
	key                  *string            // 本服务的 etcd key
	value                *ValueJson         // 本服务的 etcd value
}

// NewOptions 新的Options
func NewOptions() *options {
	return &options{}
}

func (p *options) WithAddrs(addrs []string) *options {
	p.addrs = p.addrs[0:0]
	p.addrs = append(p.addrs, addrs...)
	return p
}

func (p *options) WithTTL(ttl int64) *options {
	p.ttl = &ttl
	return p
}

func (p *options) WithGrantLeaseMaxRetries(retries int) *options {
	p.grantLeaseMaxRetries = &retries
	return p
}

func (p *options) WithDialTimeout(dialTimeout time.Duration) *options {
	p.dialTimeout = &dialTimeout
	return p
}

func (p *options) WithEventChan(eventChan chan<- interface{}) *options {
	p.eventChan = eventChan
	return p
}

func (p *options) WithWatchKeyPrefix(watchKeyPrefix string) *options {
	p.watchKeyPrefix = &watchKeyPrefix
	return p
}

func (p *options) WithKey(key string) *options {
	p.key = &key
	return p
}

func (p *options) WithValue(value *ValueJson) *options {
	p.value = value
	return p
}

func mergeOptions(opts ...*options) *options {
	no := NewOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if len(opt.addrs) != 0 {
			no.WithAddrs(opt.addrs)
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
		if opt.eventChan != nil {
			no.WithEventChan(opt.eventChan)
		}
		if opt.watchKeyPrefix != nil {
			no.WithWatchKeyPrefix(*opt.watchKeyPrefix)
		}
		if opt.key != nil {
			no.WithKey(*opt.key)
		}
		if opt.value != nil {
			no.WithValue(opt.value)
		}
	}
	return no
}

// 配置
func configure(opts *options) error {
	if len(opts.addrs) == 0 {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.ttl == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.grantLeaseMaxRetries == nil {
		var v = grantLeaseMaxRetriesDefault
		opts.grantLeaseMaxRetries = &v
	}
	if opts.dialTimeout == nil {
		opts.WithDialTimeout(dialTimeoutDefault)
	}
	if opts.eventChan == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.watchKeyPrefix == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.key == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	if opts.value == nil {
		return errors.WithMessage(xerror.Param, xruntime.Location())
	}
	return nil
}
