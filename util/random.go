package util

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	xerror "github.com/75912001/xlib/error"
	xruntime "github.com/75912001/xlib/runtime"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	mathrand "math/rand/v2"
)

// ============================================
// 全局随机数生成器（线程安全）
// ============================================
var (
	// 游戏逻辑用的快速随机数生成器（使用ChaCha8，线程安全）
	fastRand = mathrand.New(mathrand.NewChaCha8(newCryptoSeed()))
)

// 添加初始化函数
func init() {
	// 开启UUID随机（使用crypto/rand）
	uuid.EnableRandPool()
}

// newCryptoSeed 使用crypto/rand生成高质量种子
func newCryptoSeed() [32]byte {
	var seed [32]byte
	_, err := cryptorand.Read(seed[:])
	if err != nil {
		panic("failed to generate crypto seed: " + err.Error())
	}
	return seed
}

// ============================================
// 快速随机数函数（游戏逻辑用）
// 适用场景：游戏掉落、战斗计算等非安全敏感场景
// ============================================

// RandomInt 生成范围内的随机值 [min, max]
func RandomInt(min, max int) int {
	if min > max {
		min, max = max, min
	}
	return fastRand.IntN(max-min+1) + min
}

// RandomInt64 生成64位随机整数 [min, max]
func RandomInt64(min, max int64) int64 {
	if min > max {
		min, max = max, min
	}
	return fastRand.Int64N(max-min+1) + min
}

// RandomUint32 生成32位随机整数
func RandomUint32() uint32 {
	return fastRand.Uint32()
}

// RandomUint64 生成64位随机整数
func RandomUint64() uint64 {
	return fastRand.Uint64()
}

// RandomString 生成随机字符串
//
//	参数:
//		length:需要生成的长度
func RandomString(length uint32) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	const charsetLen = len(charset)
	result := make([]byte, length)
	for i := uint32(0); i < length; i++ {
		result[i] = charset[fastRand.IntN(charsetLen)]
	}
	return string(result)
}

// RandomWeighted 从权重中选出序号.[0 ... ]
//
//	参数:
//		weights:权重
//	返回值:
//		idx:weights 的序号 idx
//	e.g.: weights = [1, 2, 3], 则返回 0, 1, 2 的概率分别为 1/6, 2/6, 3/6
//	[❕] 权重为 0 的 数据 不会被选中
func RandomWeighted[T IWeight](weights []T) (idx int, err error) {
	var sum uint64
	for _, v := range weights {
		sum += uint64(v)
	}
	if sum == 0 { // weights slice 中 无数据 || 数据都为0
		return 0, errors.WithMessagef(xerror.Param, "weights sum is 0 %v", xruntime.Location())
	}
	r := fastRand.Uint64N(sum) + 1
	for i, v := range weights {
		if r <= uint64(v) {
			return i, nil
		}
		r -= uint64(v)
	}
	return 0, errors.WithMessagef(xerror.System, "random weighted error %v", xruntime.Location())
}

type IWeight interface {
	int | uint32 | uint64
}

// RandomValueBySlice 生成 随机值
//
//	参数:
//		except:排除 数据
//		slice:从该slice中随机一个,与except不重复
//	返回值:
//		slice 中的值
//	e.g.: except = [1, 2, 3], slice = [1, 2, 3, 4, 5], 则返回 4 或 5
func RandomValueBySlice(except, slice []any, equals func(a, b any) bool) any {
	var newSlice []any
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
		return nil
	}
	return newSlice[fastRand.IntN(len(newSlice))]
}

// ============================================
// 安全随机数函数（敏感操作用）
// 适用场景：密钥、Token、会话ID等安全敏感场景
// ============================================

// SecureRandomBytes 生成密码学安全的随机字节
// 适用场景：密钥、Token、会话ID等安全敏感场景
func SecureRandomBytes(length int) []byte {
	b := make([]byte, length)
	_, _ = cryptorand.Read(b)
	return b
}

// SecureRandomString 生成密码学安全的随机字符串
// 适用场景：Token、验证码等安全敏感场景
func SecureRandomString(length uint32) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	const charsetLen = len(charset)

	b := SecureRandomBytes(int(length))

	result := make([]byte, length)
	for i := uint32(0); i < length; i++ {
		result[i] = charset[int(b[i])%charsetLen]
	}
	return string(result)
}

// SecureRandomInt64 生成密码学安全的64位随机整数
func SecureRandomInt64() int64 {
	var b [8]byte
	_, _ = cryptorand.Read(b[:])
	return int64(binary.BigEndian.Uint64(b[:]))
}

// SecureRandomUint32 生成密码学安全的32位随机整数
func SecureRandomUint32() uint32 {
	var b [4]byte
	_, _ = cryptorand.Read(b[:])
	return binary.BigEndian.Uint32(b[:])
}

// SecureRandomUint64 生成密码学安全的64位随机整数
func SecureRandomUint64() uint64 {
	var b [8]byte
	_, _ = cryptorand.Read(b[:])
	return binary.BigEndian.Uint64(b[:])
}

// ============================================
// UUID 生成（已经是密码学安全的）
// ============================================

// UUIDRandomString UUID 生成 随机字符串
func UUIDRandomString() string {
	genUUID, _ := uuid.NewRandom()
	return genUUID.String()
}

// UUIDRandomBytes 生成UUID字节数组
func UUIDRandomBytes() ([16]byte, error) {
	genUUID, err := uuid.NewRandom()
	if err != nil {
		return [16]byte{}, errors.WithMessagef(err, "uuid new random %v", xruntime.Location())
	}
	var result [16]byte
	copy(result[:], genUUID[:])
	return result, nil
}
