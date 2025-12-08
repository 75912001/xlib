package selector

import (
	"context"
	xerror "github.com/75912001/xlib/error"
	xgrpcresolve "github.com/75912001/xlib/grpc/resolve"
	"github.com/75912001/xlib/grpc/util"
	xgrpcutil "github.com/75912001/xlib/grpc/util"
	xlog "github.com/75912001/xlib/log"
	xruntime "github.com/75912001/xlib/runtime"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"hash/fnv"
)

type Mod[K util.IKey] struct {
	cache *lru.Cache[CacheKey[K], util.IClientConn]
}

// 缓存淘汰回调函数
func evict[K util.IKey, V util.IClientConn](key CacheKey[K], value V) {
	xlog.PrintfInfo("grpc selector cache evict key:%v value:%v", key, value)
}

func newMod[K util.IKey]() *Mod[K] {
	cache, _ := lru.NewWithEvict[CacheKey[K], util.IClientConn](1024*1024, evict)
	return &Mod[K]{
		cache: cache,
	}
}

func (p *Mod[K]) Select(ctx context.Context, key K, method string) (*grpc.ClientConn, error) {
	ck := CacheKey[K]{
		method: method,
		key:    key,
	}
	if conn, ok := p.cache.Get(ck); ok { // 如果缓存中存在连接，则直接返回
		return conn.GetClientConn(), nil
	}
	m := xgrpcutil.Method{
		Method:      method,
		PackageName: "",
		ServiceName: "",
		MethodName:  "",
	}
	err := m.Parse()
	if err != nil {
		return nil, errors.WithMessagef(xerror.GRPCInvalidMethod, "method %s parse error: %v", method, err)
	}
	packetServiceName := "/" + m.PackageName + "." + m.ServiceName
	// 如果缓存中不存在连接，则需要从服务注册中心获取 gRPC 客户端列表
	clientConnSlice := xgrpcresolve.GetClientConn(packetServiceName)
	if len(clientConnSlice) == 0 {
		return nil, errors.WithMessage(xerror.NotExist, xruntime.Location())
	}
	idx := 0
	switch v := any(key).(type) {
	case string:
		h := fnv.New32()
		_, _ = h.Write([]byte(v))
		idx = int(h.Sum32())
	case int32:
		idx = int(v)
	case int64:
		idx = int(v)
	case uint32:
		idx = int(v)
	case uint64:
		idx = int(v)
	default:
		return nil, errors.WithMessage(xerror.NotSupport, xruntime.Location())
	}
	// 计算
	idx %= len(clientConnSlice)
	clientConn := clientConnSlice[idx]
	// 更新到缓存中
	p.cache.Add(ck, clientConn)
	return clientConn.GetClientConn(), nil
}
