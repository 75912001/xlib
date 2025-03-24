package constants

const (
	PacketLengthDefault     uint32 = 1024 * 8 // 单包最大长度 8KB
	SendChanCapacityDefault uint32 = 1024     // 发送 channel 容量-默认. 1024 大约占用0.016MB
	ServerInfoTimeOutSec    int64  = 60       // 信息-打印 超时时间 秒
	AvailableLoadDefault    uint32 = 1000000  // 可用负载-默认. 1000000
)
