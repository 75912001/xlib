package time

import (
	"strconv"
	"time"
)

// GetDayStartTimestampFromTime 当天开始时间戳
func (p *Mgr) GetDayStartTimestampFromTime(t *time.Time) int64 {
	if p.UTCSwitch.IsOn() {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).Unix()
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
}

// GetDayStartTimestampFromTimestamp 返回给定时间戳所在天的开始时间戳
func (p *Mgr) GetDayStartTimestampFromTimestamp(timestamp int64) int64 {
	if p.UTCSwitch.IsOn() {
		t := time.Unix(timestamp, 0).UTC()
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).Unix()
	}
	t := time.Unix(timestamp, 0)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
}

// GetYMDFromTimestamp 获取 e.g.:20210819
//
//	返回YMD
func (p *Mgr) GetYMDFromTimestamp(timestamp int64) int {
	var strYMD string
	if p.UTCSwitch.IsOn() {
		strYMD = time.Unix(timestamp, 0).UTC().Format("20060102")
	} else {
		strYMD = time.Unix(timestamp, 0).Format("20060102")
	}
	ymd, _ := strconv.Atoi(strYMD)
	return ymd
}

// GenYYYYMMDD 获取yyyymmdd
func GenYYYYMMDD(timestamp int64) (yyyymmdd int) {
	strYYYYMMDD := time.Unix(timestamp, 0).Format("20060102")
	yyyymmdd, _ = strconv.Atoi(strYYYYMMDD)
	return
}
