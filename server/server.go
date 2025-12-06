package server

import (
	"context"
)

type IServer interface {
	PreStart(ctx context.Context, opts ...*Options) (err error) // 服务启动前的处理 - 资源准备
	Start(ctx context.Context) (err error)                      // 启动服务
	PostStart() (err error)                                     // 服务启动后的处理 - 资源准备完成
	PreStop() (err error)                                       // 服务关闭前的处理 - 关闭资源前
	Stop() (err error)                                          // 停止服务 - 关闭 bus 之后

	GetOptions() (opt *Options)
}
