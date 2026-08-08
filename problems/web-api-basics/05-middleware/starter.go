package main

import (
	"net/http"
)

// withPoweredByHeader は、レスポンスに X-Powered-By ヘッダーを付与するミドルウェアです。
func withPoweredByHeader(next http.Handler) http.Handler {
	// TODO:
	// 1. レスポンスヘッダー "X-Powered-By" に "CodeForge" を設定する
	// 2. next.ServeHTTP(w, r) を呼び出して元の処理につなげる
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", helloHandler)

	http.ListenAndServe(":8080", withPoweredByHeader(mux))
}
