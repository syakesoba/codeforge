package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func requestUser(t *testing.T, id string) map[string]any {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+id, nil)
	rec := httptest.NewRecorder()

	newRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/v1/users/%s: 期待したステータスコードは %d でしたが、実際は %d でした（グループのパスは /api/v1 です）", id, http.StatusOK, rec.Code)
	}

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("レスポンスボディのJSONデコードに失敗しました: %v (body=%q)", err, rec.Body.String())
	}
	return got
}

func TestGetUserHandler(t *testing.T) {
	got := requestUser(t, "42")

	if got["id"] != "42" {
		t.Fatalf(`idフィールドが一致しません。期待値: "42", 実際: %v`, got["id"])
	}
	if got["name"] != "user-42" {
		t.Fatalf(`nameフィールドが一致しません。期待値: "user-42", 実際: %v`, got["name"])
	}
}

func TestGetUserHandlerDifferentID(t *testing.T) {
	got := requestUser(t, "abc")

	if got["id"] != "abc" {
		t.Fatalf(`idフィールドが一致しません。期待値: "abc", 実際: %v`, got["id"])
	}
	if got["name"] != "user-abc" {
		t.Fatalf(`nameフィールドが一致しません。期待値: "user-abc", 実際: %v`, got["name"])
	}
}
