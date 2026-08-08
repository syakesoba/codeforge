//go:build ignore

package main

import "sync"

// squareAll は nums の各要素を2乗した結果を返します。
// ゴルーチンと sync.WaitGroup を使って並行に計算してください。
func squareAll(nums []int) []int {
	result := make([]int, len(nums))

	var wg sync.WaitGroup
	for i, n := range nums {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result[i] = n * n
		}()
	}
	wg.Wait()

	return result
}

func main() {}
