package common

type ConnOptions struct {
	ReadBuffer  *int // readBuffer sets the size of the operating system's receive buffer associated with the connection. [default]: 系统默认
	WriteBuffer *int // writeBuffer sets the size of the operating system's transmit buffer associated with the connection. [default]: 系统默认
}

// NewConnOptions 新的ConnOptions
func NewConnOptions() *ConnOptions {
	return &ConnOptions{
		ReadBuffer:  nil,
		WriteBuffer: nil,
	}
}
