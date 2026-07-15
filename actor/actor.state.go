package actor

const (
	StateNew      uint32 = iota // 创建
	StateRunning                // 运行
	StateStopping               // 停止中
	StateStopped                // 已停止
)
