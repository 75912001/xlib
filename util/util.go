package util

import (
	"bytes"
	"math"
	"reflect"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"
	"unsafe"

	xerror "github.com/75912001/xlib/error"
	xpool "github.com/75912001/xlib/pool"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

// IsLittleEndian 是否小端
func IsLittleEndian() bool {
	var i int32 = 0x01020304
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb
	return b == 0x04
}

// If 三目运算符
//
//	e.g.: If(true, func() int { return 1 }, func() int { return 2 }) => 1
func If[T any](condition bool, trueFn, falseFn func() T) T {
	if condition {
		return trueFn()
	}
	return falseFn()
}

// IsDuplicate 是否有重复
//
//	e.g.: []int{1, 2, 3, 1} => true
func IsDuplicate[T comparable](slice []T) bool {
	set := make(map[T]struct{}, len(slice))
	for _, v := range slice {
		if _, exists := set[v]; exists {
			return true
		}
		set[v] = struct{}{}
	}
	return false
}

// IsDuplicateCustom 是否有重复-用于不可比较类型
//
//	[⚠️] 性能不高,慎用
//	e.g.: [1, 2, 3, 1] => true
func IsDuplicateCustom(slice []any, equals func(a, b any) bool) bool {
	set := make(map[any]struct{})
	for _, v1 := range slice {
		for v2 := range set {
			if equals(v1, v2) {
				return true
			}
		}
		set[v1] = struct{}{}
	}
	return false
}

// GetFuncName 获取函数名称
func GetFuncName(i any, seps ...rune) string {
	if i == nil {
		return xerror.Nil.Name()
	}
	funcName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	fields := strings.FieldsFunc(funcName, func(sep rune) bool {
		return slices.Contains(seps, sep)
	})
	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return xerror.Unknown.Name()
}

// PBMerge Protobuf - 深拷贝
func PBMerge(src, dst proto.Message) {
	proto.Reset(dst)
	proto.Merge(dst, src)
}

// HexStringToUint32 将十六进制字符串转换为 uint32
//
//	支持 "0x" 和 "0X" 前缀
func HexStringToUint32(hexStr string) (uint32, error) {
	if len(hexStr) > 2 && (hexStr[:2] == "0x" || hexStr[:2] == "0X") {
		hexStr = hexStr[2:]
	}
	// Parse the hex string to uint32
	value, err := strconv.ParseUint(hexStr, 16, 32)
	if err != nil {
		return 0, errors.WithMessagef(err, "hex string to uint32 %v %v", hexStr, xruntime.Location())
	}
	return uint32(value), nil
}

func PushEventWithTimeout(eventChan chan<- any, event any, timeout time.Duration) error {
	select {
	case eventChan <- event:
		return nil // 立即成功, 不分配 timer
	default:
		// 只有写不进去时才分配 timer
	}
	timer := xpool.Timer.Get()
	ok := timer.Reset(timeout)
	if !ok { // 重置失败, 重新创建一个定时器
		xpool.Timer.Put(timer)         // 旧的定时器-回收
		timer = time.NewTimer(timeout) // 新的定时器
		defer func() {
			timer.Stop()
		}()
	} else { // 重置成功, 使用已有的定时器
		defer func() {
			timer.Stop()
			xpool.Timer.Put(timer)
		}()
	}
	select {
	case eventChan <- event:
	case <-timer.C:
		return errors.WithMessagef(xerror.ChannelFull, "push event with timeout event:%v timeout:%v %v", event, timeout, xruntime.Location())
	}
	return nil
}

// TrimLeftBuffer 从左侧裁剪字节切片, 并在全部裁剪且容量过大时重新分配内存
//
//	buf: 原字节切片
//	trimLen: 需要从左侧裁剪的长度
//	maxCap: 容量阈值, 超过则重新分配
//	返回: 裁剪后的字节切片
func TrimLeftBuffer(buf []byte, trimLen, maxCap int) []byte {
	if len(buf) <= trimLen { // 全部裁剪
		if maxCap < cap(buf) { // 占用空间过大, 重新分配
			return make([]byte, 0)
		} else {
			return buf[:0]
		}
	}
	return buf[trimLen:]
}

// AdjustBufferSize 调整缓冲区大小
//
//	扩容策略:
//		如果剩余空间小于 minSpace, 则扩容
//	缩容策略:
//		如果使用的大小, 小于缓冲区大小的1/4, 并且initSize小于缓冲区大小, 则缩容为缓冲区大小的一半
//	参数:
//		buf: 原缓冲区
//		usedSize: 使用的大小
//		minSpace: 最小剩余空间要求（通常为1024）
//		initSize: 初始缓冲区大小（通常为2048）
//	返回:
//		[]byte: 调整后的缓冲区
func AdjustBufferSize(buf []byte, usedSize, minSpace, initSize int) []byte {
	remainingSpace := len(buf) - usedSize // 剩余空间
	if remainingSpace < minSpace {        // 需要扩容
		newSize := len(buf) + max(minSpace, len(buf)/2)
		newBuf := make([]byte, newSize)
		copy(newBuf, buf[:usedSize])
		return newBuf
	} else if (usedSize < len(buf)/4) && (initSize < len(buf)) { // 需要缩容
		if initSize < remainingSpace { // 确保缩容后仍有足够空间
			newBuf := make([]byte, len(buf)/2)
			copy(newBuf, buf[:usedSize])
			return newBuf
		}
	}
	return buf
}

// GenToken 生成token
//
//	参数:
//		prefix: token前缀，可为空
//	返回值:
//		token: 安全的随机字符串
func GenToken(prefix string) string {
	randomStr := SecureRandomString(32)
	if prefix != "" {
		return prefix + "." + randomStr
	}
	return randomStr
}

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

// GetGoroutineID 获取当前协程的ID
func GetGoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		return 0
	}
	b = b[:i]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

const float32Epsilon = 1e-6

// Float32Equal 判断两个 float32 是否相等
func Float32Equal(a, b float32) bool {
	return math.Abs(float64(a-b)) <= float32Epsilon
}
