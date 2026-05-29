package error

import (
	"net"

	"github.com/pkg/errors"
)

// IsNetErrorTimeout checks if a network error is a timeout.
func IsNetErrorTimeout(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

// IsNetErrClosing checks if a network error is due to a closed connection.
// 使用 errors.Is 以识别包装后的 net.ErrClosed，见 net.ErrClosed 文档说明。
func IsNetErrClosing(err error) bool {
	return errors.Is(err, net.ErrClosed)
}

func AppendError(dst, src error) error {
	if dst == nil {
		return src
	}
	if src == nil {
		return dst
	}
	return errors.Errorf("%v; %v", dst, src)
}
