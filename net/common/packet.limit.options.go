package common

type PacketLimitOptions struct {
	NewPacketLimitFunc func(maxCntPerSec uint32) IPacketLimit // 创建包限制器
	MaxCntPerSec       uint32                                 // 最大包数量-每秒
}

func NewPacketLimitOptions() *PacketLimitOptions {
	return &PacketLimitOptions{}
}

func (p *PacketLimitOptions) WithNewPacketLimitFunc(newPacketLimitFunc func(maxCntPerSec uint32) IPacketLimit) *PacketLimitOptions {
	p.NewPacketLimitFunc = newPacketLimitFunc
	return p
}

func (p *PacketLimitOptions) WithMaxCntPerSec(maxCntPerSec uint32) *PacketLimitOptions {
	p.MaxCntPerSec = maxCntPerSec
	return p
}

func (p *PacketLimitOptions) Merge(opts ...*PacketLimitOptions) *PacketLimitOptions {
	for _, opt := range opts {
		if opt.NewPacketLimitFunc != nil {
			p.NewPacketLimitFunc = opt.NewPacketLimitFunc
		}
		if opt.MaxCntPerSec != 0 {
			p.MaxCntPerSec = opt.MaxCntPerSec
		}
	}
	return p
}

func (p *PacketLimitOptions) Configure() error {
	if p.NewPacketLimitFunc != nil {
		if p.MaxCntPerSec == 0 {
			p.MaxCntPerSec = 100 // 默认100个包-每秒
		}
	}
	return nil
}
