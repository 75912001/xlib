package interceptor

import (
	"context"
	xerror "github.com/75912001/xlib/error"
	xgrpcproto "github.com/75912001/xlib/grpc/proto"
	xgrpcprotoregistry "github.com/75912001/xlib/grpc/proto/registry"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strconv"
)

func ShardKeyServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok { // 从 metadata 中获取 shareKey
			if values := md.Get(xgrpcproto.ShardKeyFieldNameDefault); len(values) > 0 {
				// 根据类型转换值
				opt := xgrpcprotoregistry.GetOptions(info.FullMethod)
				var value any
				var err error

				switch opt.ShardKeyFieldType {
				case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_STRING:
					value = values[0]
				case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_INT32:
					value, err = strconv.ParseInt(values[0], 10, 32)
					if err != nil {
						value = int32(0)
					}
				case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_INT64:
					value, err = strconv.ParseInt(values[0], 10, 64)
					if err != nil {
						value = int64(0)
					}
				case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_UINT32:
					value, err = strconv.ParseUint(values[0], 10, 32)
					if err != nil {
						value = uint32(0)
					}
				case xgrpcproto.ShardKeyFieldType_ShardKeyFieldType_UINT64:
					value, err = strconv.ParseUint(values[0], 10, 64)
					if err != nil {
						value = uint64(0)
					}
				default:
					return nil, errors.WithMessage(xerror.NotSupport, xruntime.Location())
				}
				ctx = context.WithValue(ctx, xgrpcproto.ShardKeyFieldNameDefault, value)
				return handler(ctx, req)
			}
		}
		return nil, errors.WithMessagef(xerror.GRPCNotFoundShardKey, "shard key not found for method %s", info.FullMethod)
	}
}
