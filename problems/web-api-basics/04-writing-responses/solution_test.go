package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateGreetingHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/greeting", nil)
	rec := httptest.NewRecorder()

	createGreetingHandler(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("期待したステータスコードは %d でしたが、実際は %d でした", http.StatusCreated, res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Fatalf("Content-Typeヘッダーに application/json が含まれていません。実際: %q", contentType)
	}

	var got greetingResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("レスポンスボディのJSONデコードに失敗しました: %v", err)
	}

	const want = "Hello, CodeForge!"
	if got.Message != want {
		t.Fatalf("messageフィールドが一致しません。期待値: %q, 実際: %q", want, got.Message)
	}
}
