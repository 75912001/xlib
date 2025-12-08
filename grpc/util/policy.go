package util

import (
	"context"
	"google.golang.org/grpc"
)

type IPolicy[K IKey] interface {
	Select(ctx context.Context, key K, method string) (*grpc.ClientConn, error)
}
