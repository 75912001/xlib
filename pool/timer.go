package pool

import (
	"time"
)

var Timer = NewPool(
	func() *time.Timer {
		return time.NewTimer(0)
	},
	nil,
)
