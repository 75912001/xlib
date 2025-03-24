package constants

import "time"

const (
	BusAddTimeoutDuration            = 3 * time.Second // 向总线 channel 中加入的超时时间
	BusChannelCapacityDefault uint32 = 1000000         // 总线 channel 容量-默认. 1000000 100w 大约占用15.6MB
)
