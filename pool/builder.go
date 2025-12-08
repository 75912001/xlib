package pool

import (
	"strings"
)

var Builder = NewPool(
	func() *strings.Builder {
		return &strings.Builder{}
	},
	func(buf *strings.Builder) {
		buf.Reset()
	},
)
