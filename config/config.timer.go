package config

import (
	xtimerconstants "github.com/75912001/xlib/timer/constants"
	"time"
)

type Timer struct {
	// 秒级定时器扫描间隔 [default]: timer.ScanSecondDurationDefault
	// YAML 须为 Go duration 字面量,如 100ms, 1s 不可写裸整数纳秒(yaml 无法解码为 Duration)
	ScanSecondDuration *time.Duration `yaml:"scanSecondDuration"`
	// 毫秒级定时器扫描间隔 [default]: timer.ScanMillisecondDurationDefault
	ScanMillisecondDuration *time.Duration `yaml:"scanMillisecondDuration"`
	// 毫秒级定时器实现类型 [default]: timer.MillisecondTypeList
	MillisecondType *xtimerconstants.MillisecondType `yaml:"millisecondType"`
}

func (p *Timer) GetScanSecondDuration() time.Duration {
	if p.ScanSecondDuration != nil {
		return *p.ScanSecondDuration
	}
	return xtimerconstants.ScanSecondDurationDefault
}

func (p *Timer) GetScanMillisecondDuration() time.Duration {
	if p.ScanMillisecondDuration != nil {
		return *p.ScanMillisecondDuration
	}
	return xtimerconstants.ScanMillisecondDurationDefault
}

func (p *Timer) GetMillisecondType() xtimerconstants.MillisecondType {
	if p.MillisecondType != nil {
		return *p.MillisecondType
	}
	return xtimerconstants.MillisecondTypeList
}

func (p *Timer) Configure() error {
	if p.ScanSecondDuration == nil {
		defaultValue := xtimerconstants.ScanSecondDurationDefault
		p.ScanSecondDuration = &defaultValue
	}
	if p.ScanMillisecondDuration == nil {
		defaultValue := xtimerconstants.ScanMillisecondDurationDefault
		p.ScanMillisecondDuration = &defaultValue
	}
	if p.MillisecondType == nil {
		defaultValue := xtimerconstants.MillisecondTypeList
		p.MillisecondType = &defaultValue
	}
	return nil
}
