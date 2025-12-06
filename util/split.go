package util

import (
	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

// Split2Slice 拆分字符串, 返回 ISplitValue 类型的 slice
//
//	示例:
//		Split2Slice[uint32]("1,2,3", ",")        => []uint32{1, 2, 3}
//		Split2Slice[string]("a,b,c", ",")        => []string{"a", "b", "c"}
//		Split2Slice[int]("-1,0,1", ",")          => []int{-1, 0, 1}
func Split2Slice[T ISplitValue](src, sep string) (result []T, err error) {
	if len(src) == 0 {
		return result, nil
	}

	slice := strings.Split(src, sep)
	result = make([]T, 0, len(slice))

	var t T
	switch any(t).(type) {
	case string:
		for _, v := range slice {
			result = append(result, any(v).(T))
		}
	case int:
		for _, v := range slice {
			var i64 int64
			if i64, err = strconv.ParseInt(v, 10, 64); err != nil {
				return result, errors.WithMessagef(err, "strconv parse src:%v sep:%v %v", src, sep, xruntime.Location())
			}
			result = append(result, any(int(i64)).(T))
		}
	case int32:
		for _, v := range slice {
			var i64 int64
			if i64, err = strconv.ParseInt(v, 10, 32); err != nil {
				return result, errors.WithMessagef(err, "strconv parse src:%v sep:%v %v", src, sep, xruntime.Location())
			}
			result = append(result, any(int32(i64)).(T))
		}
	case int64:
		for _, v := range slice {
			var i64 int64
			if i64, err = strconv.ParseInt(v, 10, 64); err != nil {
				return result, errors.WithMessagef(err, "strconv parse src:%v sep:%v %v", src, sep, xruntime.Location())
			}
			result = append(result, any(i64).(T))
		}
	case uint:
		for _, v := range slice {
			var u64 uint64
			if u64, err = strconv.ParseUint(v, 10, 64); err != nil {
				return result, errors.WithMessagef(err, "strconv parse src:%v sep:%v %v", src, sep, xruntime.Location())
			}
			result = append(result, any(uint(u64)).(T))
		}
	case uint32:
		for _, v := range slice {
			var u64 uint64
			if u64, err = strconv.ParseUint(v, 10, 32); err != nil {
				return result, errors.WithMessagef(err, "strconv parse src:%v sep:%v %v", src, sep, xruntime.Location())
			}
			result = append(result, any(uint32(u64)).(T))
		}
	case uint64:
		for _, v := range slice {
			var u64 uint64
			if u64, err = strconv.ParseUint(v, 10, 64); err != nil {
				return result, errors.WithMessagef(err, "strconv parse src:%v sep:%v %v", src, sep, xruntime.Location())
			}
			result = append(result, any(u64).(T))
		}
	}
	return result, nil
}

// Split2Map 拆分字符串, 返回key为 ISplitKey 类型、val为 ISplitValue 类型的map
//
//	示例:
//		Split2Map[uint32, uint32]("1,10;2,20", ";", ",")    => map[uint32]uint32{1:10, 2:20}
//		Split2Map[string, string]("k1,v1;k2,v2", ";", ",")  => map[string]string{"k1":"v1", "k2":"v2"}
//		Split2Map[string, int]("min,-100;max,100", ";", ",") => map[string]int{"min":-100, "max":100}
func Split2Map[K ISplitKey, V ISplitValue](src, sep1, sep2 string) (map[K]V, error) {
	slice := strings.Split(src, sep1)
	m := make(map[K]V, len(slice))

	var err error
	for _, v := range slice {
		if len(v) == 0 {
			continue
		}
		sliceAttr := strings.Split(v, sep2)
		if len(sliceAttr) != 2 {
			return nil, errors.WithMessagef(xerror.Param, "string split not equal 2 src:%v sep1:%v sep2:%v %v", src, sep1, sep2, xruntime.Location())
		}

		// 解析 Key
		var key K
		switch any(key).(type) {
		case string:
			key = any(sliceAttr[0]).(K)
		case int:
			var i64 int64
			if i64, err = strconv.ParseInt(sliceAttr[0], 10, 64); err != nil {
				return nil, errors.WithMessagef(err, "strconv parse key src:%v sep1:%v sep2:%v %v", src, sep1, sep2, xruntime.Location())
			}
			key = any(int(i64)).(K)
		case int32:
			var i64 int64
			if i64, err = strconv.ParseInt(sliceAttr[0], 10, 32); err != nil {
				return nil, errors.WithMessagef(err, "strconv parse key src:%v sep1:%v sep2:%v %v", src, sep1, sep2, xruntime.Location())
			}
			key = any(int32(i64)).(K)
		case int64:
			var i64 int64
			if i64, err = strconv.ParseInt(sliceAttr[0], 10, 64); err != nil {
				return nil, errors.WithMessagef(err, "strconv parse key src:%v sep1:%v sep2:%v %v", src, sep1, sep2, xruntime.Location())
			}
			key = any(i64).(K)
		case uint:
			var u64 uint64
			if u64, err = strconv.ParseUint(sliceAttr[0], 10, 64); err != nil {
				return nil, errors.WithMessagef(err, "strconv parse key src:%v sep1:%v sep2:%v %v", src, sep1, sep2, xruntime.Location())
			}
			key = any(uint(u64)).(K)
		case uint32:
			var u64 uint64
			if u64, err = strconv.ParseUint(sliceAttr[0], 10, 32); err != nil {
				return nil, errors.WithMessagef(err, "strconv parse key src:%v sep1:%v sep2:%v %v", src, sep1, sep2, xruntime.Location())
			}
			key = any(uint32(u64)).(K)
		case uint64:
			var u64 uint64
			if u64, err = strconv.ParseUint(sliceAttr[0], 10, 64); err != nil {
				return nil, errors.WithMessagef(err, "strconv parse key src:%v sep1:%v sep2:%v %v", src, sep1, sep2, xruntime.Location())
			}
			key = any(u64).(K)
		}

		// 解析 Value
		var val V
		switch any(val).(type) {
		case string:
			val = any(sliceAttr[1]).(V)
		case int:
			var i64 int64
			if i64, err = strconv.ParseInt(sliceAttr[1], 10, 64); err != nil {
				return nil, errors.WithMessagef(err, "strconv parse value src:%v sep1:%v sep2:%v %v", src, sep1, sep2, xruntime.Location())
			}
			val = any(int(i64)).(V)
		case int32:
			var i64 int64
			if i64, err = strconv.ParseInt(sliceAttr[1], 10, 32); err != nil {
				return nil, errors.WithMessagef(err, "strconv parse value src:%v sep1:%v sep2:%v %v", src, sep1, sep2, xruntime.Location())
			}
			val = any(int32(i64)).(V)
		case int64:
			var i64 int64
			if i64, err = strconv.ParseInt(sliceAttr[1], 10, 64); err != nil {
				return nil, errors.WithMessagef(err, "strconv parse value src:%v sep1:%v sep2:%v %v", src, sep1, sep2, xruntime.Location())
			}
			val = any(i64).(V)
		case uint:
			var u64 uint64
			if u64, err = strconv.ParseUint(sliceAttr[1], 10, 64); err != nil {
				return nil, errors.WithMessagef(err, "strconv parse value src:%v sep1:%v sep2:%v %v", src, sep1, sep2, xruntime.Location())
			}
			val = any(uint(u64)).(V)
		case uint32:
			var u64 uint64
			if u64, err = strconv.ParseUint(sliceAttr[1], 10, 32); err != nil {
				return nil, errors.WithMessagef(err, "strconv parse value src:%v sep1:%v sep2:%v %v", src, sep1, sep2, xruntime.Location())
			}
			val = any(uint32(u64)).(V)
		case uint64:
			var u64 uint64
			if u64, err = strconv.ParseUint(sliceAttr[1], 10, 64); err != nil {
				return nil, errors.WithMessagef(err, "strconv parse value src:%v sep1:%v sep2:%v %v", src, sep1, sep2, xruntime.Location())
			}
			val = any(u64).(V)
		}

		m[key] = val
	}
	return m, nil
}

type ISplitKey interface {
	int | int32 | int64 | uint | uint32 | uint64 | string
}

type ISplitValue interface {
	int | int32 | int64 | uint | uint32 | uint64 | string
}
