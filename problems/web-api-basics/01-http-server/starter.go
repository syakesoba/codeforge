package main

import (
	"net/http"
)

// helloHandler は "/hello" にアクセスされたときに挨拶メッセージを返すハンドラです。
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: レスポンスボディに "Hello, CodeForge!" を書き込んでください
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":8080", nil)
}
