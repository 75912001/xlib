package util

import (
	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

// Split2SliceU32 拆分字符串, 返回 uint32 类型的 slice
// e.g.: "1,2,3" => []uint32{1, 2, 3}
func Split2SliceU32(s string, sep string) (u32Slice []uint32, err error) {
	if 0 == len(s) {
		return u32Slice, nil
	}

	slice := strings.Split(s, sep)
	var u64 uint64
	for _, v := range slice {
		if u64, err = strconv.ParseUint(v, 10, 32); err != nil {
			return u32Slice, errors.WithMessage(err, xruntime.Location())
		}
		u32Slice = append(u32Slice, uint32(u64))
	}
	return u32Slice, nil
}

// Split2MapU32U32 拆分字符串, 返回key为uint32类型、val为uint32类型的map
// e.g.: "1,2;3,4" => map[uint32]int64{1:2, 3:4}
func Split2MapU32U32(s string, sep1 string, sep2 string) (map[uint32]uint32, error) {
	slice := strings.Split(s, sep1)
	m := make(map[uint32]uint32)
	var err error
	for _, v := range slice {
		if 0 == len(v) {
			continue
		}
		sliceAttr := strings.Split(v, sep2)
		if len(sliceAttr) != 2 {
			return nil, errors.WithMessage(xerror.Param, xruntime.Location())
		}
		var id uint64
		var val uint64
		if id, err = strconv.ParseUint(sliceAttr[0], 10, 32); err != nil {
			return nil, errors.WithMessage(err, xruntime.Location())
		}
		if val, err = strconv.ParseUint(sliceAttr[1], 10, 32); err != nil {
			return nil, errors.WithMessage(err, xruntime.Location())
		}
		m[uint32(id)] = uint32(val)
	}
	return m, nil
}
