//go:build ignore

package main

import (
	"context"
	"time"
)

// runWithTimeout は work を実行し、ctx がタイムアウトする前に終われば結果を、
// タイムアウトしてしまったら ctx.Err() を返します。
func runWithTimeout(ctx context.Context, work func() int) (int, error) {
	resultCh := make(chan int, 1)
	go func() {
		resultCh <- work()
	}()

	select {
	case result := <-resultCh:
		return result, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

// sumWithDeadline は nums の合計を timeout 時間以内に計算します。
func sumWithDeadline(nums []int, timeout time.Duration) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	work := func() int {
		sum := 0
		for _, n := range nums {
			sum += n
		}
		return sum
	}

	return runWithTimeout(ctx, work)
}

func main() {}
