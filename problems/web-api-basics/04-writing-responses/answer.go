//go:build ignore

package main

import (
	"encoding/json"
	"net/http"
)

// greetingResponse はJSONレスポンスの形を表す型です。
type greetingResponse struct {
	Message string `json:"message"`
}

// createGreetingHandler はJSON形式で201 Createdのレスポンスを返します。
func createGreetingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(greetingResponse{Message: "Hello, CodeForge!"})
}

func main() {
	http.HandleFunc("/greeting", createGreetingHandler)
	http.ListenAndServe(":8080", nil)
}
