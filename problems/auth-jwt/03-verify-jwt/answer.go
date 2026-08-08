//go:build ignore

package main

import (
	"errors"
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
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	return claims.GetSubject()
}

func main() {}
