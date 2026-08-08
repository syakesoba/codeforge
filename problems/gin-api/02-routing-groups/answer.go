//go:build ignore

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// getUserHandler はパスパラメータ id を読み取り、
// {"id": "<id>", "name": "user-<id>"} を200で返します。
func getUserHandler(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"id":   id,
		"name": "user-" + id,
	})
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
