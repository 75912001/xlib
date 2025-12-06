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
		strValue = (shardKeyValue).(string)
		grpcClientConn, err = stringSelectors.Sel(ctx, shardKeyValue.(string), method)
	case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_INT32:
		strValue = strconv.FormatInt(int64(shardKeyValue.(int32)), 10)
		grpcClientConn, err = int32Selectors.Sel(ctx, shardKeyValue.(int32), method)
	case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_INT64:
		strValue = strconv.FormatInt(shardKeyValue.(int64), 10)
		grpcClientConn, err = int64Selectors.Sel(ctx, shardKeyValue.(int64), method)
	case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_UINT32:
		strValue = strconv.FormatUint(uint64(shardKeyValue.(uint32)), 10)
		grpcClientConn, err = uint32Selectors.Sel(ctx, shardKeyValue.(uint32), method)
	case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_UINT64:
		strValue = strconv.FormatUint(shardKeyValue.(uint64), 10)
		grpcClientConn, err = uint64Selectors.Sel(ctx, shardKeyValue.(uint64), method)
	default:
		return ctx, nil, errors.WithMessage(xerror.GRPCNotSupportShardKeyType, xruntime.Location())
	}
	if err == nil {
		ctx = xgrpcproto.SetFromOutgoingContext(ctx, xgrpcproto.ShardKeyFieldNameDefault, strValue)
	}
	return ctx, grpcClientConn, err
}
