package timer

import "time"

var (
	ScanSecondDurationDefault      = time.Millisecond * 100 // 定时器扫描间隔-默认. 100ms
	ScanMillisecondDurationDefault = time.Millisecond * 25  // 定时器扫描间隔-默认. 25ms
)
