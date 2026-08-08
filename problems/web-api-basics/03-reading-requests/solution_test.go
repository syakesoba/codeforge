package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/add", strings.NewReader(`{"a":3,"b":4}`))
	rec := httptest.NewRecorder()

	addHandler(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("期待したステータスコードは %d でしたが、実際は %d でした", http.StatusOK, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("レスポンスボディの読み取りに失敗しました: %v", err)
	}

	const want = "7"
	if string(body) != want {
		t.Fatalf("レスポンスボディが一致しません。期待値: %q, 実際: %q", want, string(body))
	}
}

func TestAddHandlerDifferentNumbers(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/add", strings.NewReader(`{"a":10,"b":-2}`))
	rec := httptest.NewRecorder()

	addHandler(rec, req)

	body, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		t.Fatalf("レスポンスボディの読み取りに失敗しました: %v", err)
	}

	const want = "8"
	if string(body) != want {
		t.Fatalf("入力が異なる場合も正しく計算される必要があります。期待値: %q, 実際: %q", want, string(body))
	}
}
