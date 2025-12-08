package constants

const (
	ServerConfigFileSuffix = "yaml" // 服务配置文件-后缀
)

type ProcessingMode uint32 // 处理-模式

const (
	ProcessingModeBus   ProcessingMode = 0
	ProcessingModeActor ProcessingMode = 1
)
