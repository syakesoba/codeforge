package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func postUser(t *testing.T, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	newRouter().ServeHTTP(rec, req)
	return rec
}

func TestCreateUserSuccess(t *testing.T) {
	rec := postUser(t, `{"name":"Alice","email":"alice@example.com"}`)

	if rec.Code != http.StatusCreated {
		t.Fatalf("期待したステータスコードは %d でしたが、実際は %d でした (body=%q)", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("レスポンスボディのJSONデコードに失敗しました: %v (body=%q)", err, rec.Body.String())
	}
	if got["name"] != "Alice" {
		t.Fatalf(`nameフィールドが一致しません。期待値: "Alice", 実際: %v`, got["name"])
	}
	if got["email"] != "alice@example.com" {
		t.Fatalf(`emailフィールドが一致しません。期待値: "alice@example.com", 実際: %v`, got["email"])
	}
}

func TestCreateUserMissingName(t *testing.T) {
	rec := postUser(t, `{"email":"alice@example.com"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("nameが無い場合は %d を期待しますが、実際は %d でした。binding:\"required\" は付いていますか？ (body=%q)",
			http.StatusBadRequest, rec.Code, rec.Body.String())
	}
}

func TestCreateUserInvalidEmail(t *testing.T) {
	rec := postUser(t, `{"name":"Alice","email":"not-an-email"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("emailの形式が不正な場合は %d を期待しますが、実際は %d でした。binding:\"required,email\" は付いていますか？ (body=%q)",
			http.StatusBadRequest, rec.Code, rec.Body.String())
	}
}

func TestCreateUserBrokenJSON(t *testing.T) {
	rec := postUser(t, `{"name":`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("壊れたJSONの場合は %d を期待しますが、実際は %d でした (body=%q)",
			http.StatusBadRequest, rec.Code, rec.Body.String())
	}
}

func TestCreateUserErrorResponseHasErrorField(t *testing.T) {
	rec := postUser(t, `{"email":"alice@example.com"}`)

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("エラーレスポンスのJSONデコードに失敗しました: %v (body=%q)", err, rec.Body.String())
	}
	if _, ok := got["error"]; !ok {
		t.Fatalf(`エラーレスポンスには "error" フィールドが必要です。実際: %v`, got)
	}
}
