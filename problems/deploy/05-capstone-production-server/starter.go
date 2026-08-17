package main

import "net/http"

// Config はアプリケーションの設定です（Lesson 1）。
type Config struct {
	Port  string
	Debug bool
}

// HealthChecker はアプリの正常性を判定する関数の型です（Lesson 2）。
type HealthChecker func() error

// HealthHandler は実装済みです（Lesson 2の模範解答と同じ内容）。
func HealthHandler(check HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := check(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

// DebugHandler はデバッグ用の情報を返す、実装済みのハンドラーです。
// Config.Debug が true のときだけ登録します（本番で誤って有効化されないようにするため）。
func DebugHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("debug mode enabled"))
}

// NewApp は、本番運用を意識したHTTPサーバーのルーティングを組み立てます。
// これまでのレッスンで学んだ要素（設定・ヘルスチェック・DIによる差し替え可能性）を
// すべて組み合わせます。
//
//   - "/healthz" に HealthHandler(check) を登録する
//   - routes に含まれる各パス・ハンドラーをすべて登録する
//   - cfg.Debug が true の場合に限り、"/debug" に DebugHandler を登録する
//     （falseの場合は登録しない）
func NewApp(cfg Config, check HealthChecker, routes map[string]http.HandlerFunc) *http.ServeMux {
	// TODO:
	// mux := http.NewServeMux()
	// mux.HandleFunc("/healthz", HealthHandler(check))
	// routes をfor-rangeし、mux.HandleFunc(path, handler) を1つずつ登録する
	// cfg.Debug が true なら mux.HandleFunc("/debug", DebugHandler) も登録する
	// mux を返す
	return nil
}

func main() {}
