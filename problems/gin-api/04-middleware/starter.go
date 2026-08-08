package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIKeyMiddleware は X-API-Key ヘッダーを検証するミドルウェアです。
func APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO:
		// 1. c.GetHeader("X-API-Key") でヘッダーの値を取得する
		// 2. "codeforge-key" と一致しない場合は
		//    c.AbortWithStatusJSON(401, gin.H{"error": "invalid api key"}) を呼んで return する
		// 3. 一致する場合は c.Next() を呼ぶ
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
