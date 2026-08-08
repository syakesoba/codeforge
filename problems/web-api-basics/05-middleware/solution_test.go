package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithPoweredByHeader(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from inner handler"))
	})

	wrapped := withPoweredByHeader(inner)

	req := httptest.NewRequest(http.MethodGet, "/anything", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	res := rec.Result()

	got := res.Header.Get("X-Powered-By")
	if got != "CodeForge" {
		t.Fatalf("X-Powered-Byヘッダーが一致しません。期待値: %q, 実際: %q", "CodeForge", got)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("レスポンスボディの読み取りに失敗しました: %v", err)
	}

	const want = "hello from inner handler"
	if string(body) != want {
		t.Fatalf("ミドルウェアが next のハンドラを正しく呼び出していません。期待値: %q, 実際: %q", want, string(body))
	}
}
