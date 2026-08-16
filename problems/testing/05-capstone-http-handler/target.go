package main

import (
	"fmt"
	"net/http"
)

// greetHandler はクエリパラメータ "name" を読み取り、挨拶文を返します
// （実装済み・変更不要）。name が空の場合は 400 Bad Request を返します。
func greetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Hello, %s!", name)
}
