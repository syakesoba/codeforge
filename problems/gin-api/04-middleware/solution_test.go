package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getSecret(t *testing.T, apiKey string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/secret", nil)
	if apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}
	rec := httptest.NewRecorder()

	newRouter().ServeHTTP(rec, req)
	return rec
}

func TestValidAPIKeyPassesThrough(t *testing.T) {
	rec := getSecret(t, "codeforge-key")

	if rec.Code != http.StatusOK {
		t.Fatalf("正しいAPIキーの場合は %d を期待しますが、実際は %d でした。c.Next() を呼んでいますか？ (body=%q)",
			http.StatusOK, rec.Code, rec.Body.String())
	}

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("レスポンスボディのJSONデコードに失敗しました: %v (body=%q)", err, rec.Body.String())
	}
	if got["secret"] != "42" {
		t.Fatalf("secretHandlerのレスポンスが返っていません。実際: %v", got)
	}
}

func TestInvalidAPIKeyIsRejected(t *testing.T) {
	rec := getSecret(t, "wrong-key")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("誤ったAPIキーの場合は %d を期待しますが、実際は %d でした (body=%q)",
			http.StatusUnauthorized, rec.Code, rec.Body.String())
	}

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("エラーレスポンスのJSONデコードに失敗しました: %v (body=%q)", err, rec.Body.String())
	}
	if got["error"] != "invalid api key" {
		t.Fatalf(`errorフィールドが一致しません。期待値: "invalid api key", 実際: %v`, got["error"])
	}
	if _, ok := got["secret"]; ok {
		t.Fatal("認証に失敗しているのに secretHandler が実行されています。c.Abort系のメソッドで中断していますか？")
	}
}

func TestMissingAPIKeyIsRejected(t *testing.T) {
	rec := getSecret(t, "")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("APIキーが無い場合は %d を期待しますが、実際は %d でした (body=%q)",
			http.StatusUnauthorized, rec.Code, rec.Body.String())
	}
}
