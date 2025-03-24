package timer

import (
	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"time"
)

// NewOptions 新的Options
func NewOptions() *options {
	return &options{}
}

// options contains options to configure instance. Each option can be set through setter functions. See
// documentation for each setter Function for an explanation of the option.
type options struct {
	scanSecondDuration      *time.Duration     // 扫描秒级定时器,纳秒间隔(如 100000000,则每100毫秒扫描一次秒定时器)
	scanMillisecondDuration *time.Duration     // 扫描毫秒级定时器,纳秒间隔(如 100000000,则每100毫秒扫描一次毫秒定时器)
	outgoingTimeoutChan     chan<- interface{} // 是超时事件放置的channel,由外部传入.超时的 Second/millisecond 都会放入其中
}

func (p *options) WithScanSecondDuration(scanSecondDuration time.Duration) *options {
	p.scanSecondDuration = &scanSecondDuration
	return p
}

func (p *options) WithScanMillisecondDuration(scanMillisecondDuration time.Duration) *options {
	p.scanMillisecondDuration = &scanMillisecondDuration
	return p
}

func (p *options) WithOutgoingTimerOutChan(timeoutChan chan<- interface{}) *options {
	p.outgoingTimeoutChan = timeoutChan
	return p
}

// merge combines the given *options into a single *options in a last one wins fashion.
// The specified options are merged with the existing options on the server, with the specified options taking
// precedence.
func (p *options) merge(opts ...*options) *options {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.scanSecondDuration != nil {
			p.scanSecondDuration = opt.scanSecondDuration
		}
		if opt.scanMillisecondDuration != nil {
			p.scanMillisecondDuration = opt.scanMillisecondDuration
		}
		if opt.outgoingTimeoutChan != nil {
			p.outgoingTimeoutChan = opt.outgoingTimeoutChan
		}
	}
	return p
}

// 配置
func (p *options) configure() error {
	if p.scanSecondDuration == nil {
		p.scanSecondDuration = &ScanSecondDurationDefault
	}
	if p.scanMillisecondDuration == nil {
		p.scanMillisecondDuration = &ScanMillisecondDurationDefault
	}
	if p.outgoingTimeoutChan == nil {
		return errors.WithMessagef(xerror.ChannelNil, xruntime.Location())
	}
	return nil
}
