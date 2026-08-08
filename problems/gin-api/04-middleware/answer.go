//go:build ignore

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIKeyMiddleware は X-API-Key ヘッダーを検証するミドルウェアです。
func APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key != "codeforge-key" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
			return
		}
		c.Next()
	}
}

func secretHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"secret": "42"})
}

// newRouter はこのレッスンで使うルーティングを登録したルーターを返します。
func newRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/secret", APIKeyMiddleware(), secretHandler)

	return r
}

func main() {
	newRouter().Run(":8080")
}
