package main

import (
	"net/http"
)

// greetingResponse はJSONレスポンスの形を表す型です。
type greetingResponse struct {
	Message string `json:"message"`
}

// createGreetingHandler はJSON形式で201 Createdのレスポンスを返します。
func createGreetingHandler(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Content-Type ヘッダーに "application/json" を設定する
	// 2. ステータスコード 201 (http.StatusCreated) を書き込む
	// 3. greetingResponse{Message: "Hello, CodeForge!"} をJSONとして書き込む
	// ※ 順番は ヘッダー設定 → ステータスコード → 本文 です
}

func main() {
	http.HandleFunc("/greeting", createGreetingHandler)
	http.ListenAndServe(":8080", nil)
}
