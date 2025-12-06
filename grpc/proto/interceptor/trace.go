package interceptor

import (
	"context"
	xgrpcproto "github.com/75912001/xlib/grpc/proto"
	xutil "github.com/75912001/xlib/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TraceServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok { // 从 metadata 中获取 traceID
			if values := md.Get(xgrpcproto.TraceIdFieldNameDefault); len(values) > 0 {
				ctx = context.WithValue(ctx, xgrpcproto.TraceIdFieldNameDefault, values[0])
				return handler(ctx, req)
			}
		}
		// 如果没有找到 traceID，则生成一个新的
		traceID := xutil.UUIDRandomString()
		ctx = context.WithValue(ctx, xgrpcproto.TraceIdFieldNameDefault, traceID)
		return handler(ctx, req)
	}
}

func TraceClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var traceID string
		if val, ok := ctx.Value(xgrpcproto.TraceIdFieldNameDefault).(string); ok { // 获取 traceID
			traceID = val
		} else { // 生成新的 traceID
			traceID = xutil.UUIDRandomString()
			ctx = context.WithValue(ctx, xgrpcproto.TraceIdFieldNameDefault, traceID)
		}
		// 获取现有的 metadata
		md, _ := metadata.FromOutgoingContext(ctx)
		if md == nil {
			md = metadata.New(nil)
		}
		md.Set(xgrpcproto.TraceIdFieldNameDefault, traceID)
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
