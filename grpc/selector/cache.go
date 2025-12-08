package selector

import (
	"github.com/75912001/xlib/grpc/util"
)

// CacheKey 缓存键
type CacheKey[K util.IKey] struct {
	method string
	key    K
}
