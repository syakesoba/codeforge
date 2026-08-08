package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGreetHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/greet/Gopher", nil)
	rec := httptest.NewRecorder()

	newMux().ServeHTTP(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("期待したステータスコードは %d でしたが、実際は %d でした", http.StatusOK, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("レスポンスボディの読み取りに失敗しました: %v", err)
	}

	const want = "Hello, Gopher!"
	if string(body) != want {
		t.Fatalf("レスポンスボディが一致しません。期待値: %q, 実際: %q", want, string(body))
	}
}

func TestGreetHandlerDifferentName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/greet/Codeforge", nil)
	rec := httptest.NewRecorder()

	newMux().ServeHTTP(rec, req)

	body, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		t.Fatalf("レスポンスボディの読み取りに失敗しました: %v", err)
	}

	const want = "Hello, Codeforge!"
	if string(body) != want {
		t.Fatalf("パスパラメータが異なる場合も正しく反映される必要があります。期待値: %q, 実際: %q", want, string(body))
	}
}
