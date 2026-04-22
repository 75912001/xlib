package interceptor

import (
	"context"
	xgrpcproto "github.com/75912001/xlib/grpc/proto"
	xgrpcprotoregistry "github.com/75912001/xlib/grpc/proto/registry"
	"google.golang.org/grpc"
	"time"
)

func TimeOutClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 获取超时时间
		opt := xgrpcprotoregistry.GetOptions(method)
		duration, err := time.ParseDuration(opt.Timeout)
		if err != nil || duration <= 0 {
			// 非法或空配置时避免 WithTimeout(..., 0) 导致立刻失败，与 GetOptions 默认语义对齐
			duration = xgrpcproto.RpcTimeoutDurationDefault
		}
		ctx, cancel := context.WithTimeout(ctx, duration)
		defer func() {
			cancel()
		}()
		// 调用下一个拦截器或实际的 RPC 调用
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
