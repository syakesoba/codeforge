package main

import (
	"context"
	"time"
)

// Shutdowner は正常終了処理を行うインターフェースです。
// *http.Server はこのインターフェースを満たす Shutdown(ctx) メソッドを持っています。
type Shutdowner interface {
	Shutdown(ctx context.Context) error
}

// GracefulShutdown は、timeout以内に shutdowner の終了処理が完了することを試みます。
//
//   - timeout以内に Shutdown が完了すれば、その戻り値（nilまたはエラー）をそのまま返す
//   - timeoutを超過した場合、Shutdown自体はcontextのキャンセルを検知してエラーを返すので、
//     その結果をそのまま返せばよい（GracefulShutdown自身がタイムアウトを別途判定する必要はない）
func GracefulShutdown(shutdowner Shutdowner, timeout time.Duration) error {
	// TODO:
	// ctx, cancel := context.WithTimeout(context.Background(), timeout) でタイムアウト付きのcontextを作る
	// defer cancel() を忘れずに
	// return shutdowner.Shutdown(ctx)
	return nil
}

func main() {}
