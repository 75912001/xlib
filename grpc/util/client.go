package util

import (
	"google.golang.org/grpc"
)

type IClientConn interface {
	GetClientConn() *grpc.ClientConn // 获取 gRPC 客户端连接
	Disabled()
	Available() bool // 检查服务是否可用
	Stop() error     // 停止客户端连接，释放资源
	GetID() string
}
