package main

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret-key"

func TestIssueTokenReturnsJWTFormat(t *testing.T) {
	tokenString, err := IssueToken(testSecret, "alice", time.Hour)
	if err != nil {
		t.Fatalf("IssueToken がエラーを返しました: %v", err)
	}
	if tokenString == "" {
		t.Fatal("IssueToken が空文字列を返しました")
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		t.Fatalf("JWTはドット区切りの3パートであるべきですが、%d パートでした: %q", len(parts), tokenString)
	}
}

func TestIssueTokenIsVerifiableWithSecret(t *testing.T) {
	tokenString, err := IssueToken(testSecret, "alice", time.Hour)
	if err != nil {
		t.Fatalf("IssueToken がエラーを返しました: %v", err)
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(testSecret), nil
	})
	if err != nil {
		t.Fatalf("発行されたトークンを秘密鍵で検証できません: %v", err)
	}
	if !token.Valid {
		t.Fatal("発行されたトークンが無効と判定されました")
	}
}

func TestIssueTokenHasSubjectClaim(t *testing.T) {
	tokenString, err := IssueToken(testSecret, "bob", time.Hour)
	if err != nil {
		t.Fatalf("IssueToken がエラーを返しました: %v", err)
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(testSecret), nil
	})
	if err != nil {
		t.Fatalf("トークンの検証に失敗しました: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("クレームを jwt.MapClaims として取得できません")
	}

	sub, err := claims.GetSubject()
	if err != nil {
		t.Fatalf(`"sub" クレームが取得できません: %v`, err)
	}
	if sub != "bob" {
		t.Fatalf(`"sub" クレームが一致しません。期待値: "bob", 実際: %q`, sub)
	}
}

func TestIssueTokenHasExpiryClaim(t *testing.T) {
	tokenString, err := IssueToken(testSecret, "alice", time.Hour)
	if err != nil {
		t.Fatalf("IssueToken がエラーを返しました: %v", err)
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(testSecret), nil
	})
	if err != nil {
		t.Fatalf("トークンの検証に失敗しました: %v", err)
	}

	claims := token.Claims.(jwt.MapClaims)
	exp, err := claims.GetExpirationTime()
	if err != nil {
		t.Fatalf(`"exp" クレームが取得できません: %v`, err)
	}
	if exp == nil {
		t.Fatal(`"exp" クレームが設定されていません。有効期限は必ず設定してください`)
	}

	// 1時間後前後（誤差5分以内）であること
	want := time.Now().Add(time.Hour)
	diff := exp.Time.Sub(want)
	if diff > 5*time.Minute || diff < -5*time.Minute {
		t.Fatalf("有効期限が想定と異なります。期待値: 約 %v, 実際: %v", want, exp.Time)
	}
}

func TestIssueTokenExpiredWhenNegativeTTL(t *testing.T) {
	tokenString, err := IssueToken(testSecret, "alice", -time.Hour)
	if err != nil {
		t.Fatalf("IssueToken がエラーを返しました: %v", err)
	}

	// 期限切れのトークンは検証に失敗するはず
	_, err = jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(testSecret), nil
	})
	if err == nil {
		t.Fatal("期限切れのトークンが有効と判定されました。exp クレームを設定していますか？")
	}
}

func TestIssueTokenRejectedWithWrongSecret(t *testing.T) {
	tokenString, err := IssueToken(testSecret, "alice", time.Hour)
	if err != nil {
		t.Fatalf("IssueToken がエラーを返しました: %v", err)
	}

	_, err = jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte("wrong-secret"), nil
	})
	if err == nil {
		t.Fatal("誤った秘密鍵で検証が成功してしまいました。secret で署名していますか？")
	}
}
