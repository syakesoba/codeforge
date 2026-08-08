//go:build ignore

package main

// collectSquares は nums の各要素を2乗した値をchannel経由で送信し、
// 呼び出し側でrangeを使って集めた結果を返します。
func collectSquares(nums []int) []int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for _, n := range nums {
			ch <- n * n
		}
	}()

	var result []int
	for v := range ch {
		result = append(result, v)
	}
	return result
}

func main() {}
