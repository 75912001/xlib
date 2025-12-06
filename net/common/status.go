package common

import (
	xcontrol "github.com/75912001/xlib/control"
	"time"
)

// 链接状态

// 活跃 -> 非活跃
//			非活跃 -> 死亡
//			非活跃 -> 活跃

type Status struct {
	Inactive       bool               // 非活跃状态
	InactiveStart  int64              // 非活跃开始时间戳
	InactiveData   []byte             // 非活跃期间缓存的数据
	DeathTimestamp int64              // 死亡时间戳
	DeathCallback  xcontrol.ICallBack // 死亡时回调函数
}

func NewStatus() *Status {
	return &Status{
		Inactive:       false,
		InactiveStart:  0,
		InactiveData:   nil,
		DeathTimestamp: 0,
		DeathCallback:  nil,
	}
}

// 设置活跃
func (p *Status) SetActive() {
	p.Inactive = false
	p.InactiveStart = 0
	p.InactiveData = nil
	p.DeathTimestamp = 0
	p.DeathCallback = nil
}

// 设置非活跃
func (p *Status) SetInactive(deathDurationSecond int64, deathCallback xcontrol.ICallBack) {
	p.Inactive = true
	p.InactiveStart = time.Now().Unix()
	p.InactiveData = make([]byte, 0)
	p.DeathTimestamp = p.InactiveStart + deathDurationSecond
	p.DeathCallback = deathCallback
}

// append 缓存
func (p *Status) AppendCache(data []byte) {
	p.InactiveData = append(p.InactiveData, data...)
}

// 获取缓存
func (p *Status) GetCache() []byte {
	return p.InactiveData
}
