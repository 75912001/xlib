package pool

import (
	"time"
)

var Timer = NewPool(
	func() *time.Timer {
		return time.NewTimer(0)
	},
	func(t *time.Timer) {
		if !t.Stop() {
			select {
			case <-t.C:
			default:
			}
		}
	},
)
