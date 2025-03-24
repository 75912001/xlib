package timer

import "time"

var timestampOffset int64 // 时间戳偏移量

// SetTimestampOffset 设置时间戳偏移量
func SetTimestampOffset(offset int64) {
	timestampOffset = offset
}

// ShadowTimestamp 叠加偏移量的时间戳-秒
func ShadowTimestamp() int64 {
	return time.Now().Unix() + timestampOffset
}
