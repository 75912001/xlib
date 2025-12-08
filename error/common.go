package error

import (
	"net"
	"strings"
)

// IsNetErrorTimeout checks if a network error is a timeout.
func IsNetErrorTimeout(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

// IsNetErrClosing checks if a network error is due to a closed connection.
func IsNetErrClosing(err error) bool {
	return err != nil && strings.Contains(err.Error(), "use of closed network connection")
}
