package server

import "context"

type IServer interface {
	Start(ctx context.Context, opts ...*ServerOptions) (err error) // 启动服务
	PreStop() (err error)                                          // 服务关闭前的处理 - 关闭资源前
	Stop() (err error)                                             // 停止服务 - 关闭 bus 之后
}
