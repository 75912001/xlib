package server

import "context"

type IServer interface {
	Start(ctx context.Context) (err error) // 启动服务
	PreStop() (err error)                  // 服务关闭前的处理
	Stop() (err error)                     // 停止服务
}
