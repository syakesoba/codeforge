package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const testSecret = "test-secret-key"

func getMe(t *testing.T, authHeader string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	rec := httptest.NewRecorder()

	newRouter(testSecret).ServeHTTP(rec, req)
	return rec
}

func TestValidTokenIsAccepted(t *testing.T) {
	token, err := IssueToken(testSecret, "alice", time.Hour)
	if err != nil {
		t.Fatalf("IssueToken がエラーを返しました: %v", err)
	}

	rec := getMe(t, "Bearer "+token)

	if rec.Code != http.StatusOK {
		t.Fatalf("正しいトークンの場合は %d を期待しますが、実際は %d でした。c.Next() を呼んでいますか？ (body=%q)",
			http.StatusOK, rec.Code, rec.Body.String())
	}

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("レスポンスのJSONデコードに失敗しました: %v (body=%q)", err, rec.Body.String())
	}
	if got["username"] != "alice" {
		t.Fatalf(`usernameが一致しません。期待値: "alice", 実際: %v。c.Set("username", ...) を呼んでいますか？`, got["username"])
	}
}

func TestMissingHeaderIsRejected(t *testing.T) {
	rec := getMe(t, "")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("Authorizationヘッダーが無い場合は %d を期待しますが、実際は %d でした",
			http.StatusUnauthorized, rec.Code)
	}
}

func TestMissingBearerPrefixIsRejected(t *testing.T) {
	token, _ := IssueToken(testSecret, "alice", time.Hour)

	// "Bearer " を付けずにトークンだけ送る
	rec := getMe(t, token)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(`"Bearer " プレフィックスが無い場合は %d を期待しますが、実際は %d でした`,
			http.StatusUnauthorized, rec.Code)
	}
}

func TestInvalidTokenIsRejected(t *testing.T) {
	rec := getMe(t, "Bearer not-a-valid-token")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("不正なトークンの場合は %d を期待しますが、実際は %d でした",
			http.StatusUnauthorized, rec.Code)
	}
}

func TestExpiredTokenIsRejected(t *testing.T) {
	token, err := IssueToken(testSecret, "alice", -time.Hour)
	if err != nil {
		t.Fatalf("IssueToken がエラーを返しました: %v", err)
	}

	rec := getMe(t, "Bearer "+token)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("期限切れトークンの場合は %d を期待しますが、実際は %d でした",
			http.StatusUnauthorized, rec.Code)
	}
}

func TestTokenSignedWithWrongSecretIsRejected(t *testing.T) {
	token, err := IssueToken("attacker-secret", "attacker", time.Hour)
	if err != nil {
		t.Fatalf("IssueToken がエラーを返しました: %v", err)
	}

	rec := getMe(t, "Bearer "+token)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("異なる秘密鍵で署名されたトークンは %d を期待しますが、実際は %d でした",
			http.StatusUnauthorized, rec.Code)
	}
}

func TestRejectedRequestDoesNotReachHandler(t *testing.T) {
	rec := getMe(t, "Bearer invalid")

	var got map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("レスポンスのJSONデコードに失敗しました: %v (body=%q)", err, rec.Body.String())
	}
	if _, ok := got["username"]; ok {
		t.Fatal("認証に失敗しているのに meHandler が実行されています。c.Abort系のメソッドで中断していますか？")
	}
	if got["error"] != "unauthorized" {
		t.Fatalf(`errorフィールドが一致しません。期待値: "unauthorized", 実際: %v`, got["error"])
	}
}
