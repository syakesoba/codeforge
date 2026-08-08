package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	newRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("期待したステータスコードは %d でしたが、実際は %d でした", http.StatusOK, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Fatalf("Content-Typeヘッダーに application/json が含まれていません。実際: %q", contentType)
	}

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("レスポンスボディのJSONデコードに失敗しました: %v (body=%q)", err, rec.Body.String())
	}

	if got["status"] != "ok" {
		t.Fatalf(`statusフィールドが一致しません。期待値: "ok", 実際: %v`, got["status"])
	}
}
