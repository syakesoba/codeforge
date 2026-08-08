package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// IssueToken は username を主体とする、ttl 時間有効なJWTを発行します（実装済み）。
func IssueToken(secret, username string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

// VerifyToken はトークンを検証し、sub クレーム（ユーザー名）を返します。
func VerifyToken(secret, tokenString string) (string, error) {
	// TODO:
	// 1. jwt.Parse(tokenString, キー関数) で検証する
	//    キー関数の中で t.Method が *jwt.SigningMethodHMAC であることを確認し、
	//    そうでなければエラーを返す
	// 2. jwt.Parse がエラーを返したら、それをそのまま返す
	// 3. token.Claims を jwt.MapClaims に型アサーションする
	// 4. claims.GetSubject() の結果を返す
	return "", nil
}

func main() {}
