// 时间
// 程序运行过程中,会使用时间.计算时间.

package time

import (
	xcontrol "github.com/75912001/xlib/control"
	"time"
)

// Mgr 时间管理器
type Mgr struct {
	timestampSecond       int64                  // 上一次调用Update更新的时间戳-秒
	timestampMillisecond  int64                  // 上一次调用Update更新的时间戳-毫秒
	time                  time.Time              // 上一次调用Update更新的时间
	timestampSecondOffset int64                  // 时间戳偏移量-秒
	UTCSwitch             xcontrol.ISwitchButton // UTC 时间开关
}

func NewMgr() *Mgr {
	return &Mgr{
		UTCSwitch: xcontrol.NewSwitchButton(false),
	}
}

// NowTime 获取当前时间
func (p *Mgr) NowTime() time.Time {
	if p.UTCSwitch.IsOn() {
		return time.Now().UTC()
	}
	return time.Now()
}

// Update 更新时间管理器中的,当前时间
func (p *Mgr) Update() {
	p.time = p.NowTime()
	p.timestampSecond = p.time.Unix()
	p.timestampMillisecond = p.time.UnixMilli()
}

// ShadowTimestamp 叠加偏移量的时间戳-秒
func (p *Mgr) ShadowTimestamp() int64 {
	return p.timestampSecond + p.timestampSecondOffset
}

// SetTimestampOffset 设置 时间戳偏移量-秒
func (p *Mgr) SetTimestampOffset(offset int64) {
	p.timestampSecondOffset = offset
}
