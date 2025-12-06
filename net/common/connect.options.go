package common

type ConnOptions struct {
	ReadBuffer  *int // readBuffer sets the size of the operating system's receive buffer associated with the connection. [default]: 1024*1024/系统默认
	WriteBuffer *int // writeBuffer sets the size of the operating system's transmit buffer associated with the connection. [default]: 1024*1024/系统默认
}

// NewConnOptions 新的ConnOptions
func NewConnOptions() *ConnOptions {
	return &ConnOptions{}
}

func (p *ConnOptions) WithReadBuffer(readBuffer int) *ConnOptions {
	p.ReadBuffer = &readBuffer
	return p
}

func (p *ConnOptions) WithWriteBuffer(writeBuffer int) *ConnOptions {
	p.WriteBuffer = &writeBuffer
	return p
}

func (p *ConnOptions) Merge(opts ...*ConnOptions) *ConnOptions {
	for _, opt := range opts {
		if opt.ReadBuffer != nil {
			p.ReadBuffer = opt.ReadBuffer
		}
		if opt.WriteBuffer != nil {
			p.WriteBuffer = opt.WriteBuffer
		}
	}
	return p
}

func (p *ConnOptions) Configure() error {
	if p.ReadBuffer == nil {
		p.ReadBuffer = new(int)
		*p.ReadBuffer = 1024 * 1024
	}
	if p.WriteBuffer == nil {
		p.WriteBuffer = new(int)
		*p.WriteBuffer = 1024 * 1024
	}
	return nil
}
