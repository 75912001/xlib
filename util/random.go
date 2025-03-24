package util

import (
	"bytes"
	cryptorand "crypto/rand"
	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/pkg/errors"
	"math/big"
	"math/rand"
)

// RandomInt 生成范围内的随机值
//
//	参数:
//		min:最小值
//		max:最大值
func RandomInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

// RandomString 生成随机字符串
//
//	参数:
//		len:需要生成的长度
func RandomString(len uint32) (container string, err error) {
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	bigInt := big.NewInt(int64(bytes.NewBufferString(str).Len()))
	for i := uint32(0); i < len; i++ {
		if randomInt, err := cryptorand.Int(cryptorand.Reader, bigInt); err != nil {
			return "", errors.WithMessage(err, xruntime.Location())
		} else {
			container += string(str[randomInt.Int64()])
		}
	}
	return container, nil
}

// RandomWeighted 从权重中选出序号.[0 ... ]
//
//	[NOTE] 参数 weights 必须有长度
//	参数:
//		weights:权重
//	返回值:
//		idx:weights 的序号 idx
//
// e.g.: weights = [1, 2, 3], 则返回 0, 1, 2 的概率分别为 1/6, 2/6, 3/6
func RandomWeighted(weights []uint32) (idx int, err error) {
	var sum int64
	for _, v := range weights {
		sum += int64(v)
	}
	if sum == 0 { //weights slice 中 无数据 / 数据都为0
		return 0, errors.WithMessage(xerror.Param, xruntime.Location())
	}
	r := rand.Int63n(sum) + 1
	for i, v := range weights {
		if r <= int64(v) {
			return i, nil
		}
		r -= int64(v)
	}
	return 0, errors.WithMessage(xerror.System, xruntime.Location())
}

// RandomValueBySlice 生成 随机值
//
//	参数:
//		except:排除 数据
//		slice:从该slice中随机一个,与except不重复
//	返回值:
//		slice 中的值
//
// e.g.: except = [1, 2, 3], slice = [1, 2, 3, 4, 5], 则返回 4 或 5
func RandomValueBySlice(except []interface{}, slice []interface{}, equals func(a, b interface{}) bool) (interface{}, error) {
	var newSlice []interface{}
	for _, v := range slice {
		found := false
		for _, e := range except {
			if equals(v, e) {
				found = true
				break
			}
		}
		if !found {
			newSlice = append(newSlice, v)
		}
	}
	if len(newSlice) == 0 {
		return nil, errors.WithMessage(xerror.NotExist, xruntime.Location())
	}
	return newSlice[rand.Intn(len(newSlice))], nil
}
