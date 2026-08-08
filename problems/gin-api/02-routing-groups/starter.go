package main

import (
	"github.com/gin-gonic/gin"
)

// getUserHandler はパスパラメータ id を読み取り、
// {"id": "<id>", "name": "user-<id>"} を200で返します。
func getUserHandler(c *gin.Context) {
	// TODO:
	// 1. c.Param("id") でパスパラメータを取得する
	// 2. c.JSON でステータスコード200とJSONを返す
}

// newRouter はこのレッスンで使うルーティングを登録したルーターを返します。
func newRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/users/:id", getUserHandler)
	}

	return r
}

func main() {
	newRouter().Run(":8080")
}
