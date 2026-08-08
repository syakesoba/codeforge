//go:build ignore

package main

import (
	"net/http"
)

// withPoweredByHeader は、レスポンスに X-Powered-By ヘッダーを付与するミドルウェアです。
func withPoweredByHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "CodeForge")
		next.ServeHTTP(w, r)
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
