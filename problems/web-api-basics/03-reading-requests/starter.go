package main

import (
	"net/http"
)

// addRequest はリクエストボディ {"a": 3, "b": 4} をデコードするための型です。
type addRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

// addHandler はリクエストボディの a と b を足し算し、結果を文字列で返します。
func addHandler(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. json.NewDecoder(r.Body).Decode(...) でリクエストボディを addRequest にデコードする
	// 2. req.A + req.B の結果を文字列としてレスポンスボディに書き込む
}

func main() {
	http.HandleFunc("/add", addHandler)
	http.ListenAndServe(":8080", nil)
}
