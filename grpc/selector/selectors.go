package selector

import (
	"context"

	xerror "github.com/75912001/xlib/error"
	xgrpcproto "github.com/75912001/xlib/grpc/proto"
	xgrpcprotoregistry "github.com/75912001/xlib/grpc/proto/registry"
	xgrpcutil "github.com/75912001/xlib/grpc/util"
	xmap "github.com/75912001/xlib/map"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type selectors[K xgrpcutil.IKey] struct {
	MapMgr *xmap.MapMgr[string, xgrpcutil.IPolicy[K]] // key: /${packageName}.${serviceName}/${methodName}, value: 负载均衡策略
}

func newSelectors[K xgrpcutil.IKey]() *selectors[K] {
	var zero K
	switch any(zero).(type) {
	case string:
	case int32:
	case int64:
	case uint32:
	case uint64:
	default:
		panic(errors.WithMessagef(xerror.NotSupport, "key type %T not support", zero))
	}
	return &selectors[K]{
		MapMgr: xmap.NewMapMgr[string, xgrpcutil.IPolicy[K]](),
	}
}

var (
	stringSelectors = newSelectors[string]()
	int32Selectors  = newSelectors[int32]()
	int64Selectors  = newSelectors[int64]()
	uint32Selectors = newSelectors[uint32]()
	uint64Selectors = newSelectors[uint64]()
)

// Init initializes the selectors for different key types.
func Init() {
	stringMod := newMod[string]()
	int32Mod := newMod[int32]()
	int64Mod := newMod[int64]()
	uint32Mod := newMod[uint32]()
	uint64Mod := newMod[uint64]()

	strHashRing := newHashRing[string]()
	int32HashRing := newHashRing[int32]()
	int64HashRing := newHashRing[int64]()
	uint32HashRing := newHashRing[uint32]()
	uint64HashRing := newHashRing[uint64]()

	for k, v := range xgrpcprotoregistry.GMethodOptions {
		switch v.LoadBalancePolicy {
		case xgrpcproto.LoadBalancePolicy_LoadBalancePolicy_Mod:
			switch v.ShardKeyFieldType {
			case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_STRING:
				stringSelectors.MapMgr.Add(k, stringMod)
			case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_INT32:
				int32Selectors.MapMgr.Add(k, int32Mod)
			case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_INT64:
				int64Selectors.MapMgr.Add(k, int64Mod)
			case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_UINT32:
				uint32Selectors.MapMgr.Add(k, uint32Mod)
			case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_UINT64:
				uint64Selectors.MapMgr.Add(k, uint64Mod)
			default:
				panic(errors.WithMessagef(xerror.NotSupport, "shard key type %s not support for method %s", v.ShardKeyFieldType, k))
			}
		case xgrpcproto.LoadBalancePolicy_LoadBalancePolicy_RingHash:
			switch v.ShardKeyFieldType {
			case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_STRING:
				stringSelectors.MapMgr.Add(k, strHashRing)
			case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_INT32:
				int32Selectors.MapMgr.Add(k, int32HashRing)
			case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_INT64:
				int64Selectors.MapMgr.Add(k, int64HashRing)
			case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_UINT32:
				uint32Selectors.MapMgr.Add(k, uint32HashRing)
			case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_UINT64:
				uint64Selectors.MapMgr.Add(k, uint64HashRing)
			default:
				panic(errors.WithMessagef(xerror.NotSupport, "shard key type %s not support for method %s", v.ShardKeyFieldType, k))
			}
		default:
			panic(errors.WithMessagef(xerror.NotSupport, "load balance type %s not support for method %s", v, k))
		}
	}
}

func (p *selectors[K]) Sel(ctx context.Context, k K, method string) (*grpc.ClientConn, error) {
	policy, ok := p.MapMgr.Find(method)
	if !ok {
		return nil, errors.WithMessagef(xerror.NotFound, "selector for method %s", method)
	}
	return policy.Select(ctx, k, method)
}
