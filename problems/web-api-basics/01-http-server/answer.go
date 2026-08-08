//go:build ignore

package main

import (
	"fmt"
	"net/http"
)

// helloHandler は "/hello" にアクセスされたときに挨拶メッセージを返すハンドラです。
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, CodeForge!")
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":8080", nil)
}
