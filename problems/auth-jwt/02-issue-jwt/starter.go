package main

import (
	"time"
)

// IssueToken は username を主体とする、ttl 時間有効なJWTを発行します。
func IssueToken(secret, username string, ttl time.Duration) (string, error) {
	// TODO:
	// 1. jwt.MapClaims で "sub"（username）と "exp"（time.Now().Add(ttl).Unix()）を設定する
	// 2. jwt.NewWithClaims(jwt.SigningMethodHS256, claims) でトークンを作る
	// 3. token.SignedString([]byte(secret)) の結果を返す
	return "", nil
}

func main() {}
