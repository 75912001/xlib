package common

import "github.com/xtaci/kcp-go/v5"

type KCPOptions struct {
	SndWindowSize *int           // 窗口大小-发送 [default:512]
	RcvWindowSize *int           // 窗口大小-接收 [default:512]
	Nodelay       *int           // 无延迟 0:关闭 1:打开 [default:1]
	Interval      *int           // 间隔 [default:20ms]
	Resend        *int           // 重传模式 [default:2]
	Nc            *int           // 关闭Nagle算法 0:打开 1:关闭 [default:1]
	AckNodelay    *bool          // 关闭延迟确认 [default:true]
	Mtu           *int           // 最大传输单元 [default:1350]
	BlockCrypt    kcp.BlockCrypt //加密,解密 [default:nil]
	Fec           *bool          // 是否开启FEC [default:false]
	DataShards    *int           // Fec: true 数据分片数 [default: 10] Fec: false 数据分片数 [default: 0]
	ParityShards  *int           // Fec: true 奇偶校验分片数 [default: 3] Fec: false 奇偶校验分片数 [default: 0]
}

func (p *KCPOptions) WithSndWindowSize(sndWindowSize int) *KCPOptions {
	p.SndWindowSize = &sndWindowSize
	return p
}
func (p *KCPOptions) WithRcvWindowSize(rcvWindowSize int) *KCPOptions {
	p.RcvWindowSize = &rcvWindowSize
	return p
}
func (p *KCPOptions) WithNoDelay(nodelay int) *KCPOptions {
	p.Nodelay = &nodelay
	return p
}
func (p *KCPOptions) WithInterval(interval int) *KCPOptions {
	p.Interval = &interval
	return p
}
func (p *KCPOptions) WithResend(resend int) *KCPOptions {
	p.Resend = &resend
	return p
}
func (p *KCPOptions) WithNc(nc int) *KCPOptions {
	p.Nc = &nc
	return p
}
func (p *KCPOptions) WithAckNodelay(ackNodelay bool) *KCPOptions {
	p.AckNodelay = &ackNodelay
	return p
}
func (p *KCPOptions) WithMtu(mtu int) *KCPOptions {
	p.Mtu = &mtu
	return p
}
func (p *KCPOptions) WithBlockCrypt(blockCrypt kcp.BlockCrypt) *KCPOptions {
	p.BlockCrypt = blockCrypt
	return p
}
func (p *KCPOptions) WithFEC(fec bool) *KCPOptions {
	p.Fec = &fec
	return p
}

func (p *KCPOptions) Merge(opts ...*KCPOptions) *KCPOptions {
	for _, opt := range opts {
		if opt.SndWindowSize != nil {
			p.WithSndWindowSize(*opt.SndWindowSize)
		}
		if opt.RcvWindowSize != nil {
			p.WithRcvWindowSize(*opt.RcvWindowSize)
		}
		if opt.Nodelay != nil {
			p.WithNoDelay(*opt.Nodelay)
		}
		if opt.Interval != nil {
			p.WithInterval(*opt.Interval)
		}
		if opt.Resend != nil {
			p.WithResend(*opt.Resend)
		}
		if opt.Nc != nil {
			p.WithNc(*opt.Nc)
		}
		if opt.AckNodelay != nil {
			p.WithAckNodelay(*opt.AckNodelay)
		}
		if opt.Mtu != nil {
			p.WithMtu(*opt.Mtu)
		}
		if opt.BlockCrypt != nil {
			p.WithBlockCrypt(opt.BlockCrypt)
		}
		if opt.Fec != nil {
			p.WithFEC(*opt.Fec)
		}
	}
	return p
}

func (p *KCPOptions) Configure() error {
	if p.SndWindowSize == nil {
		p.SndWindowSize = new(int)
		*p.SndWindowSize = 512
	}
	if p.RcvWindowSize == nil {
		p.RcvWindowSize = new(int)
		*p.RcvWindowSize = 512
	}
	if p.Nodelay == nil {
		p.Nodelay = new(int)
		*p.Nodelay = 1
	}
	if p.Interval == nil {
		p.Interval = new(int)
		*p.Interval = 10
	}
	if p.Resend == nil {
		p.Resend = new(int)
		*p.Resend = 2
	}
	if p.Nc == nil {
		p.Nc = new(int)
		*p.Nc = 1
	}
	if p.AckNodelay == nil {
		p.AckNodelay = new(bool)
		*p.AckNodelay = true
	}
	if p.Mtu == nil {
		p.Mtu = new(int)
		*p.Mtu = 1350
	}
	if p.BlockCrypt == nil {
		p.BlockCrypt = nil
	}
	if p.Fec == nil {
		p.Fec = new(bool)
		*p.Fec = false
	}
	if *p.Fec {
		p.DataShards = new(int)
		*p.DataShards = 10
		p.ParityShards = new(int)
		*p.ParityShards = 3
	} else {
		p.DataShards = new(int)
		*p.DataShards = 0
		p.ParityShards = new(int)
		*p.ParityShards = 0
	}
	return nil
}
