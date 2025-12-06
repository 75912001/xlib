package timer

import (
	"math"
)

// 时间轮数量
const cycleSize int = 33

// 时间轮持续时间
//
//	key:序号[0,...]
//	value:到期时间
var cycleDuration [cycleSize]int64 //nolint:all // 包内 全局变量

func init() { //nolint:all // 初始化
	for i := range cycleDuration {
		cycleDuration[i] = genDuration(i)
	}
	cycleDuration[cycleSize-1] = math.MaxInt64
}

// 生成一个轮的时长
//
//	参数:
//		轮序号
//		返回值: 2,4,6,8,12,16,20,24,32,40,48,56,72,88,104,120,152,184,216,248,312,376,440,504,632,760,888,1016,1272,1528,1784,2040,math.MaxInt64
func genDuration(idx int) int64 {
	base := int64(2) // 初始步长
	stepCount := 4   // 步长递增的槽数间隔
	// 计算当前步长
	curStep := base << uint(idx/stepCount)
	// 计算当前槽的起始值
	prevSum := int64(0)
	for i := range idx {
		prevStep := base << uint(i/stepCount)
		prevSum += prevStep
	}
	return prevSum + curStep
}

// 根据 时长 找到时间轮的序号 二分查找 (迭代)
func searchCycleIdx(duration int64) int {
	const h = len(cycleDuration) - 1
	low, high := 0, h
	for low <= high {
		mid := low + (high-low)/2 //nolint:all // 二分法,2:从中间取
		if low == high {
			if cycleDuration[mid] < duration {
				return mid + 1
			} else {
				return mid
			}
		}
		switch {
		case cycleDuration[mid] == duration:
			return mid
		case duration < cycleDuration[mid]:
			high = mid - 1
		case duration > cycleDuration[mid]:
			low = mid + 1
		}
	}
	return low
}

// 向前查找符合时间差的时间轮序号
//
//	参数:
//		duration: 到期 时长
//		idx: 轮序号 0 < idx
func findPrevCycleIdx(duration int64, idx int) int {
	for {
		if idx != 0 && duration <= cycleDuration[idx-1] {
			idx--
		} else {
			break
		}
	}
	return idx
}
