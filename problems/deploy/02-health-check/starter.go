package main

import "net/http"

// HealthChecker は、アプリが正常に動作しているかどうかを判定する関数の型です。
// 例えばDBへの接続確認など、実際のチェック内容は呼び出し側が自由に決められます。
type HealthChecker func() error

// HealthHandler は、HealthChecker の結果に応じてレスポンスを返す
// net/http ハンドラーを組み立てて返します。
//
//   - check() がエラーを返さなければ、ステータスコード200と本文"ok"を返す
//   - check() がエラーを返せば、ステータスコード503（Service Unavailable）と
//     そのエラーメッセージ（err.Error()）を本文として返す
func HealthHandler(check HealthChecker) http.HandlerFunc {
	// TODO:
	// return func(w http.ResponseWriter, r *http.Request) { ... } を返す
	// 中身で check() を呼び、結果に応じて w.WriteHeader と w.Write（または fmt.Fprint）を使い分ける
	return nil
}

func main() {}
