package interceptor

import (
	"context"
	xgrpcprotoregistry "github.com/75912001/xlib/grpc/proto/registry"
	"google.golang.org/grpc"
	"time"
)

func TimeOutClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 获取超时时间
		opt := xgrpcprotoregistry.GetOptions(method)
		// 创建带超时的 context
		duration, _ := time.ParseDuration(opt.Timeout)
		ctx, cancel := context.WithTimeout(ctx, duration)
		defer func() {
			cancel()
		}()
		// 调用下一个拦截器或实际的 RPC 调用
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
