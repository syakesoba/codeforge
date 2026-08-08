//go:build ignore

package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

// VerifyToken はトークンを検証し、sub クレームを返します（実装済み）。
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

// AuthMiddleware は Authorization: Bearer <token> を検証するミドルウェアです。
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		username, err := VerifyToken(secret, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set("username", username)
		c.Next()
	}
}

// meHandler はコンテキストに保存されたユーザー名を返します（実装済み）。
func meHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "username not found in context"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": username})
}

// newRouter はこのレッスンで使うルーティングを登録したルーターを返します。
func newRouter(secret string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/me", AuthMiddleware(secret), meHandler)

	return r
}

func main() {
	newRouter("secret").Run(":8080")
}
