package config

import (
	xtimerconstants "github.com/75912001/xlib/timer/constants"
	"time"
)

type Timer struct {
	// 秒级定时器 扫描间隔(纳秒) 1000*1000*100=100000000 为100毫秒 [default]: timer.ScanSecondDurationDefault
	scanSecondDuration *time.Duration `yaml:"scanSecondDuration"`
	// 毫秒级定时器 扫描间隔(纳秒) 1000*1000*100=100000000 为25毫秒 [default]: timer.ScanMillisecondDurationDefault
	scanMillisecondDuration *time.Duration `yaml:"scanMillisecondDuration"`
	// 毫秒级定时器-使用类型
	millisecondType *xtimerconstants.MillisecondType `yaml:"millisecondType"` // [default]: timer.MillisecondTypeList
}

func (p *Timer) GetScanSecondDuration() time.Duration {
	if p.scanSecondDuration != nil {
		return *p.scanSecondDuration
	}
	return xtimerconstants.ScanSecondDurationDefault
}

func (p *Timer) GetScanMillisecondDuration() time.Duration {
	if p.scanMillisecondDuration != nil {
		return *p.scanMillisecondDuration
	}
	return xtimerconstants.ScanMillisecondDurationDefault
}

func (p *Timer) GetMillisecondType() xtimerconstants.MillisecondType {
	if p.millisecondType != nil {
		return *p.millisecondType
	}
	return xtimerconstants.MillisecondTypeList
}

func (p *Timer) Configure() error {
	if p.scanSecondDuration == nil {
		defaultValue := xtimerconstants.ScanSecondDurationDefault
		p.scanSecondDuration = &defaultValue
	}
	if p.scanMillisecondDuration == nil {
		defaultValue := xtimerconstants.ScanMillisecondDurationDefault
		p.scanMillisecondDuration = &defaultValue
	}
	if p.millisecondType == nil {
		defaultValue := xtimerconstants.MillisecondTypeList
		p.millisecondType = &defaultValue
	}
	return nil
}
