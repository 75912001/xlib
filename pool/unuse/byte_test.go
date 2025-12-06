package unuse

import (
	"strconv"
	"testing"
)

func BenchmarkBytePool(b *testing.B) {
	// 测试不同大小的内存分配
	//sizes := []int{64, 1024, 2048, 4096, 65536}
	sizes := []int{32, 512, 2048, 8192, 32768, 131072, 524288, 2097152}
	for _, size := range sizes {
		b.Run("WithPool-"+strconv.Itoa(size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				buf := MakeByteSlice(size)
				ReleaseByteSlice(buf)
			}
		})

		b.Run("WithoutPool-"+strconv.Itoa(size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = make([]byte, size)
			}
		})
	}
}

// 测试并发场景下的性能
func BenchmarkBytePoolConcurrent(b *testing.B) {
	size := 1024
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := MakeByteSlice(size)
			ReleaseByteSlice(buf)
		}
	})
}

// 测试边界情况
func TestBytePoolEdgeCases(t *testing.T) {
	testCases := []struct {
		name string
		size int
		want bool // true 表示期望返回非nil
	}{
		{"零大小", 0, false},
		{"正常大小", 1024, true},
		{"最大边界", maxAreaValue, true},
		{"超出边界", maxAreaValue + 1, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := MakeByteSlice(tc.size)
			if (buf != nil) != tc.want {
				t.Errorf("MakeByteSlice(%d) = %v, want %v", tc.size, buf != nil, tc.want)
			}
		})
	}
}
