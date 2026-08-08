//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// addRequest はリクエストボディ {"a": 3, "b": 4} をデコードするための型です。
type addRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

// addHandler はリクエストボディの a と b を足し算し、結果を文字列で返します。
func addHandler(w http.ResponseWriter, r *http.Request) {
	var req addRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "%d", req.A+req.B)
}

func main() {
	http.HandleFunc("/add", addHandler)
	http.ListenAndServe(":8080", nil)
}
