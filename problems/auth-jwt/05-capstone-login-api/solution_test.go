package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func doJSON(t *testing.T, r http.Handler, method, path, body, authHeader string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func decodeBody(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("レスポンスのJSONデコードに失敗しました: %v (body=%q)", err, rec.Body.String())
	}
	return got
}

func TestAuthFlow(t *testing.T) {
	r := newRouter()

	// 1. 登録できること
	{
		rec := doJSON(t, r, http.MethodPost, "/register", `{"username":"alice","password":"password123"}`, "")
		if rec.Code != http.StatusCreated {
			t.Fatalf("POST /register: 期待したステータスコードは %d でしたが、実際は %d でした (body=%q)",
				http.StatusCreated, rec.Code, rec.Body.String())
		}
		got := decodeBody(t, rec)
		if got["username"] != "alice" {
			t.Fatalf(`POST /register: usernameが一致しません。期待値: "alice", 実際: %v`, got["username"])
		}
	}

	// 2. パスワードが平文で保存されていないこと
	{
		mu.Lock()
		user := users["alice"]
		mu.Unlock()

		if user.PasswordHash == "password123" {
			t.Fatal("パスワードが平文で保存されています。HashPassword を使ってください")
		}
		if user.PasswordHash == "" {
			t.Fatal("PasswordHash が保存されていません")
		}
	}

	// 3. 同じユーザー名の重複登録は409であること
	{
		rec := doJSON(t, r, http.MethodPost, "/register", `{"username":"alice","password":"another-password"}`, "")
		if rec.Code != http.StatusConflict {
			t.Fatalf("重複登録: 期待したステータスコードは %d でしたが、実際は %d でした",
				http.StatusConflict, rec.Code)
		}
	}

	// 4. 短すぎるパスワードは400であること
	{
		rec := doJSON(t, r, http.MethodPost, "/register", `{"username":"bob","password":"short"}`, "")
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("8文字未満のパスワード: 期待したステータスコードは %d でしたが、実際は %d でした。binding:\"required,min=8\" は付いていますか？",
				http.StatusBadRequest, rec.Code)
		}
	}

	// 5. ログインしてトークンを取得できること
	var token string
	{
		rec := doJSON(t, r, http.MethodPost, "/login", `{"username":"alice","password":"password123"}`, "")
		if rec.Code != http.StatusOK {
			t.Fatalf("POST /login: 期待したステータスコードは %d でしたが、実際は %d でした (body=%q)",
				http.StatusOK, rec.Code, rec.Body.String())
		}
		got := decodeBody(t, rec)
		tokenValue, ok := got["token"].(string)
		if !ok || tokenValue == "" {
			t.Fatalf(`POST /login: "token" フィールドが返っていません。実際: %v`, got)
		}
		token = tokenValue
	}

	// 6. 誤ったパスワードは401であること
	{
		rec := doJSON(t, r, http.MethodPost, "/login", `{"username":"alice","password":"wrong-password"}`, "")
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("誤ったパスワード: 期待したステータスコードは %d でしたが、実際は %d でした",
				http.StatusUnauthorized, rec.Code)
		}
	}

	// 7. 存在しないユーザーも同じく401であること
	{
		rec := doJSON(t, r, http.MethodPost, "/login", `{"username":"nobody","password":"password123"}`, "")
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("存在しないユーザー: 期待したステータスコードは %d でしたが、実際は %d でした",
				http.StatusUnauthorized, rec.Code)
		}
	}

	// 8. トークンを使って /profile にアクセスできること
	{
		rec := doJSON(t, r, http.MethodGet, "/profile", "", "Bearer "+token)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /profile: 期待したステータスコードは %d でしたが、実際は %d でした (body=%q)",
				http.StatusOK, rec.Code, rec.Body.String())
		}
		got := decodeBody(t, rec)
		if got["username"] != "alice" {
			t.Fatalf(`GET /profile: usernameが一致しません。期待値: "alice", 実際: %v`, got["username"])
		}
	}

	// 9. 認証なしの /profile は401であること
	{
		rec := doJSON(t, r, http.MethodGet, "/profile", "", "")
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("認証なしの /profile: 期待したステータスコードは %d でしたが、実際は %d でした",
				http.StatusUnauthorized, rec.Code)
		}
	}
}
