package selector

import (
	"context"
	xerror "github.com/75912001/xlib/error"
	xgrpcproto "github.com/75912001/xlib/grpc/proto"
	xgrpcprotoregistry "github.com/75912001/xlib/grpc/proto/registry"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"strconv"
)

func Sel(ctx context.Context, method string, shardKeyValue any) (context.Context, *grpc.ClientConn, error) {
	var err error
	var grpcClientConn *grpc.ClientConn
	opt := xgrpcprotoregistry.GetOptions(method)
	switch opt.LoadBalancePolicy {
	case xgrpcproto.LoadBalancePolicy_LoadBalancePolicy_Mod:
	case xgrpcproto.LoadBalancePolicy_LoadBalancePolicy_RingHash:
	default:
		return ctx, nil, errors.WithMessagef(xerror.NotSupport, "load balance type %s not support", opt.LoadBalancePolicy)
	}
	var strValue string
	switch opt.ShardKeyFieldType {
	case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_STRING:
		v, ok := shardKeyValue.(string)
		if !ok {
			return ctx, nil, errors.WithMessagef(xerror.GRPCNotSupportShardKeyType,
				"method %s shard key: want string, got %T", method, shardKeyValue)
		}
		strValue = v
		grpcClientConn, err = stringSelectors.Sel(ctx, v, method)
	case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_INT32:
		v, ok := shardKeyValue.(int32)
		if !ok {
			return ctx, nil, errors.WithMessagef(xerror.GRPCNotSupportShardKeyType,
				"method %s shard key: want int32, got %T", method, shardKeyValue)
		}
		strValue = strconv.FormatInt(int64(v), 10)
		grpcClientConn, err = int32Selectors.Sel(ctx, v, method)
	case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_INT64:
		v, ok := shardKeyValue.(int64)
		if !ok {
			return ctx, nil, errors.WithMessagef(xerror.GRPCNotSupportShardKeyType,
				"method %s shard key: want int64, got %T", method, shardKeyValue)
		}
		strValue = strconv.FormatInt(v, 10)
		grpcClientConn, err = int64Selectors.Sel(ctx, v, method)
	case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_UINT32:
		v, ok := shardKeyValue.(uint32)
		if !ok {
			return ctx, nil, errors.WithMessagef(xerror.GRPCNotSupportShardKeyType,
				"method %s shard key: want uint32, got %T", method, shardKeyValue)
		}
		strValue = strconv.FormatUint(uint64(v), 10)
		grpcClientConn, err = uint32Selectors.Sel(ctx, v, method)
	case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_UINT64:
		v, ok := shardKeyValue.(uint64)
		if !ok {
			return ctx, nil, errors.WithMessagef(xerror.GRPCNotSupportShardKeyType,
				"method %s shard key: want uint64, got %T", method, shardKeyValue)
		}
		strValue = strconv.FormatUint(v, 10)
		grpcClientConn, err = uint64Selectors.Sel(ctx, v, method)
	default:
		return ctx, nil, errors.WithMessage(xerror.GRPCNotSupportShardKeyType, xruntime.Location())
	}
	if err == nil {
		ctx = xgrpcproto.SetFromOutgoingContext(ctx, xgrpcproto.ShardKeyFieldNameDefault, strValue)
	}
	return ctx, grpcClientConn, err
}
