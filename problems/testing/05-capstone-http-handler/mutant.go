//go:build ignore

package main

import (
	"fmt"
	"net/http"
)

// greetHandler はわざとバグを仕込んだ実装です
// （nameが空かどうかのチェックを忘れており、常に200を返してしまう）。
func greetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Hello, %s!", name)
}
