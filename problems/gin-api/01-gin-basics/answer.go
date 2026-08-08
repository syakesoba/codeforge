//go:build ignore

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// healthHandler は {"status": "ok"} をステータスコード200で返します。
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// newRouter はこのレッスンで使うルーティングを登録したルーターを返します。
func newRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/health", healthHandler)

	return r
}

func main() {
	newRouter().Run(":8080")
}
