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
	// TODO:
	// 1. context.WithTimeout(context.Background(), timeout) でctxとcancelを作る（defer cancel() を忘れずに）
	// 2. nums の合計を計算するクロージャを work として用意する
	// 3. runWithTimeout(ctx, work) を呼び出し、その結果をそのまま return する
	return 0, nil
}

func main() {}
