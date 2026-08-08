package main

import (
	"github.com/gin-gonic/gin"
)

// healthHandler は {"status": "ok"} をステータスコード200で返します。
func healthHandler(c *gin.Context) {
	// TODO: c.JSON を使って {"status": "ok"} を200で返してください
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
