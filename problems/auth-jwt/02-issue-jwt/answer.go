//go:build ignore

package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// IssueToken は username を主体とする、ttl 時間有効なJWTを発行します。
func IssueToken(secret, username string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func main() {}
