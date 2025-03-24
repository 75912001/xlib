package timer

import (
	"math"
)

// 时间轮数量
const cycleSize int = 9

// 时间轮持续时间
//
//	key:序号[0,...]
//	value:到期时间
var cycleDuration [cycleSize]int64 //nolint:all // 包内 全局变量

func init() { //nolint:all // 初始化
	for i := 0; i < len(cycleDuration); i++ {
		cycleDuration[i] = genDuration(i)
	}
	cycleDuration[cycleSize-1] = math.MaxInt64
}

// 生成一个轮的时长
//
//	参数:
//		轮序号
//	返回值:
//		4,8,16,32,64,128,256,512,math.MaxInt64
func genDuration(idx int) int64 {
	const shift = 2 // 偏移量
	return int64(1 << uint(idx+shift))
}

//// 根据 时长 找到时间轮的序号
//// (当前为从头依次判断,适用于大多数数据 符合头部条件,若数据均匀分布,则适用于使用二分查找)
////
////	参数:
////		duration:时长
////	返回值:
////		轮序号
//func findCycleIdx(duration int64) (idx int) {
//	for k, v := range cycleDuration {
//		if duration <= v {
//			return k
//		} else {
//			idx++
//		}
//	}
//	return len(cycleDuration) - 1
//}

// 根据 时长 找到时间轮的序号 二分查找 (迭代)
func searchCycleIdxIteration(duration int64) int {
	low, high := 0, len(cycleDuration)-1
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
