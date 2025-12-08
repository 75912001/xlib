package constants

import "time"

var (
	ScanSecondDurationDefault      = time.Millisecond * 100 // 定时器扫描间隔-默认. 100ms
	ScanMillisecondDurationDefault = time.Millisecond * 20  // 定时器扫描间隔-默认. 20ms
)

type MillisecondType uint32 // 毫秒级定时器使用类型

const (
	MillisecondTypeList    MillisecondType = 0 // 毫秒级定时器 使用 list.List 实现. 尾部插入-O(1) 插入排序-O(n) 头部移除-O(1)
	MillisecondTypeMinHeap MillisecondType = 1 // 毫秒级定时器 使用小顶堆实现. Push-O(log n), Pop-O(log n)
)
