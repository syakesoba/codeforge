package main

import (
	"testing"
)

func makeRange(n int) []int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i
	}
	return nums
}

func TestBuildCubeMap(t *testing.T) {
	nums := []int{1, 2, 3, -2, 0}

	got := buildCubeMap(nums)

	if len(got) != len(nums) {
		t.Fatalf("mapの要素数が一致しません。期待値: %d, 実際: %d", len(nums), len(got))
	}
	for _, n := range nums {
		want := n * n * n
		if got[n] != want {
			t.Fatalf("buildCubeMap[%d] = %d, want %d", n, got[n], want)
		}
	}
}

// 十分な数のゴルーチンから同時にmapへ書き込ませることで、
// Mutexで保護されていない実装は "fatal error: concurrent map writes" で
// プロセスごと異常終了する（go testの終了コードが非0になり不合格になる）。
func TestBuildCubeMapConcurrentSafety(t *testing.T) {
	nums := makeRange(2000)

	for i := 0; i < 20; i++ {
		got := buildCubeMap(nums)
		if len(got) != len(nums) {
			t.Fatalf("試行%d回目: mapの要素数が一致しません。期待値: %d, 実際: %d", i+1, len(nums), len(got))
		}
	}
}
