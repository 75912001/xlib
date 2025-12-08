package common

import (
	"time"
)

type IPacketLimit interface {
	IsOverload(cnt uint32, nowTime time.Time) bool // 是否超载, cnt 为当前包数量, nowTime 为当前时间
}

type PackLimitDefault struct {
	MaxCntPerSec uint32    // 最大包数量-每秒
	Cnt          uint32    // 包数量
	Time         time.Time // 时间
}

func NewPackLimitDefault(maxCntPerSec uint32) IPacketLimit {
	return &PackLimitDefault{
		MaxCntPerSec: maxCntPerSec,
		Cnt:          0,
		Time:         time.Now(),
	}
}

func (p *PackLimitDefault) IsOverload(cnt uint32, nowTime time.Time) bool {
	if p.Time.Unix() == nowTime.Unix() {
		p.Cnt += cnt
	} else {
		p.Cnt = cnt
		p.Time = nowTime
	}
	return p.MaxCntPerSec < p.Cnt
}
