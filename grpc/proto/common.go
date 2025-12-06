package proto

import (
	"context"
	"google.golang.org/grpc/metadata"
	"time"
)

var (
	RpcTimeoutDurationDefault = 5 * time.Second // 默认rpc超时时间为5秒
	RpcTimeoutDefault         = "5s"            // 默认rpc超时时间字符串表示
	ShardKeyFieldNameDefault  = "x-shard-key"
	TraceIdFieldNameDefault   = "x-trace-id"
)

func SetTraceIdFieldNameDefault(name string) {
	TraceIdFieldNameDefault = name
}

func SetRpcTimeoutDurationDefault(duration time.Duration) {
	RpcTimeoutDurationDefault = duration
}

func SetShardKeyFieldNameDefault(name string) {
	ShardKeyFieldNameDefault = name
}

func SetFromOutgoingContext(ctx context.Context, k string, val string) context.Context {
	md, _ := metadata.FromOutgoingContext(ctx)
	if md == nil {
		md = metadata.New(nil)
	}
	md.Set(k, val)
	return metadata.NewOutgoingContext(ctx, md)
}
