package main

import (
	"encoding/base64"
	"strconv"
	"strings"
	"testing"
	"time"
)

const testSecret = "test-secret-key"

func encodeSegment(s string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(s))
}

func itoa(n int64) string {
	return strconv.FormatInt(n, 10)
}

func TestVerifyTokenReturnsSubject(t *testing.T) {
	tokenString, err := IssueToken(testSecret, "alice", time.Hour)
	if err != nil {
		t.Fatalf("IssueToken がエラーを返しました: %v", err)
	}

	sub, err := VerifyToken(testSecret, tokenString)
	if err != nil {
		t.Fatalf("VerifyToken がエラーを返しました: %v", err)
	}
	if sub != "alice" {
		t.Fatalf(`subが一致しません。期待値: "alice", 実際: %q`, sub)
	}
}

func TestVerifyTokenDifferentUser(t *testing.T) {
	tokenString, err := IssueToken(testSecret, "bob", time.Hour)
	if err != nil {
		t.Fatalf("IssueToken がエラーを返しました: %v", err)
	}

	sub, err := VerifyToken(testSecret, tokenString)
	if err != nil {
		t.Fatalf("VerifyToken がエラーを返しました: %v", err)
	}
	if sub != "bob" {
		t.Fatalf(`subが一致しません。期待値: "bob", 実際: %q`, sub)
	}
}

func TestVerifyTokenRejectsWrongSecret(t *testing.T) {
	tokenString, err := IssueToken(testSecret, "alice", time.Hour)
	if err != nil {
		t.Fatalf("IssueToken がエラーを返しました: %v", err)
	}

	if _, err := VerifyToken("wrong-secret", tokenString); err == nil {
		t.Fatal("誤った秘密鍵での検証が成功してしまいました")
	}
}

func TestVerifyTokenRejectsExpired(t *testing.T) {
	// 1時間前に期限切れになったトークン
	tokenString, err := IssueToken(testSecret, "alice", -time.Hour)
	if err != nil {
		t.Fatalf("IssueToken がエラーを返しました: %v", err)
	}

	if _, err := VerifyToken(testSecret, tokenString); err == nil {
		t.Fatal("期限切れトークンの検証が成功してしまいました")
	}
}

func TestVerifyTokenRejectsGarbage(t *testing.T) {
	if _, err := VerifyToken(testSecret, "not-a-jwt-at-all"); err == nil {
		t.Fatal("JWTでない文字列の検証が成功してしまいました")
	}
}

func TestVerifyTokenRejectsTamperedToken(t *testing.T) {
	tokenString, err := IssueToken(testSecret, "alice", time.Hour)
	if err != nil {
		t.Fatalf("IssueToken がエラーを返しました: %v", err)
	}

	// ペイロード（"sub"を別のユーザーに）だけを差し替え、署名はそのままにする。
	// 攻撃者が他人になりすまそうとするケースを再現している。
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		t.Fatalf("IssueToken が返すトークンがJWT形式ではありません: %q", tokenString)
	}
	exp := time.Now().Add(time.Hour).Unix()
	parts[1] = encodeSegment(`{"sub":"attacker","exp":` + itoa(exp) + `}`)
	tampered := strings.Join(parts, ".")

	if _, err := VerifyToken(testSecret, tampered); err == nil {
		t.Fatal("改ざんされたトークンの検証が成功してしまいました")
	}
}

// alg: none 攻撃を模したトークンを拒否できることを確認する。
func TestVerifyTokenRejectsNoneAlgorithm(t *testing.T) {
	// {"alg":"none","typ":"JWT"}.{"sub":"attacker","exp":<未来>}. （署名なし）
	exp := time.Now().Add(time.Hour).Unix()
	header := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0"
	payload := encodeSegment(`{"sub":"attacker","exp":` + itoa(exp) + `}`)
	noneToken := header + "." + payload + "."

	if _, err := VerifyToken(testSecret, noneToken); err == nil {
		t.Fatal("alg:none のトークンが受け入れられました。キー関数で署名アルゴリズムを確認していますか？")
	}
}
