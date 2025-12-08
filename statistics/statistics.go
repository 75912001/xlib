package statistics

import "time"

type IStatistics interface {
	GetCount() uint64
	GetErrorCount() uint64
	GetProcessTime() time.Duration
}

type Statistics struct {
	Count       uint64        // 处理总数
	ErrorCount  uint64        // 错误总数
	ProcessTime time.Duration // 处理时间-总时间-毫秒
}

func NewStatistics() *Statistics {
	return &Statistics{}
}

func (p *Statistics) GetCount() uint64 {
	return p.Count
}
func (p *Statistics) GetErrorCount() uint64 {
	return p.ErrorCount
}
func (p *Statistics) GetProcessTime() time.Duration {
	return p.ProcessTime
}
