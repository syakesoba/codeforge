package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()

	helloHandler(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("期待したステータスコードは %d でしたが、実際は %d でした", http.StatusOK, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("レスポンスボディの読み取りに失敗しました: %v", err)
	}

	const want = "Hello, CodeForge!"
	if string(body) != want {
		t.Fatalf("レスポンスボディが一致しません。期待値: %q, 実際: %q", want, string(body))
	}
}
