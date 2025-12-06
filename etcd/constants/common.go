package constants

import "time"

var (
	GrantLeaseMaxRetriesDefault = 600             // 授权租约 最大 重试次数
	DialTimeoutDefault          = time.Second * 5 // dialTimeout is the timeout for failing to establish a connection.
)
