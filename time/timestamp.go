package time

import "strconv"

// FormatTimestampMs 将毫秒时间戳格式化为十进制字符串。
func FormatTimestampMs(timestampMs int64) string {
	return strconv.FormatInt(timestampMs, 10)
}

// ParseTimestampMs 将十进制字符串解析为毫秒时间戳。
func ParseTimestampMs(value string) (int64, bool) {
	timestampMs, err := strconv.ParseInt(value, 10, 64)
	return timestampMs, err == nil
}
