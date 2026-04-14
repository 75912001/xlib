package pool

import (
	"time"
)

var Timer = NewPool(
	func() *time.Timer {
		t := time.NewTimer(0)
		t.Stop()
		return t
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
