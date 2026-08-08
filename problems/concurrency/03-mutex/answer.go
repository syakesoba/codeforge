//go:build ignore

package main

import "sync"

// buildCubeMap は nums の各要素をキーとして、その3乗の値をmapに格納して返します。
// 複数のゴルーチンから安全に書き込めるよう、sync.Mutexで保護してください。
func buildCubeMap(nums []int) map[int]int {
	result := make(map[int]int)
	var mu sync.Mutex

	var wg sync.WaitGroup
	for _, n := range nums {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cube := n * n * n

			mu.Lock()
			result[n] = cube
			mu.Unlock()
		}()
	}
	wg.Wait()

	return result
}

func main() {}
