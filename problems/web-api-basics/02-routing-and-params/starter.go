package main

import (
	"net/http"
)

// greetHandler はパスパラメータ "name" を読み取り、挨拶を返します。
func greetHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: パスパラメータ "name" を取得し、"Hello, <name>!" と書き込んでください
}

// newMux はこのレッスンで使うルーティングを登録したServeMuxを返します。
func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /greet/{name}", greetHandler)
	return mux
}

func main() {
	http.ListenAndServe(":8080", newMux())
}
