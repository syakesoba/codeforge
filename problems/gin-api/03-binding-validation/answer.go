//go:build ignore

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// createUserRequest はリクエストボディ {"name": "...", "email": "..."} を受け取る型です。
type createUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// createUserHandler はリクエストボディをバインド・バリデーションし、
// 失敗なら400、成功なら201でレスポンスを返します。
func createUserHandler(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"name":  req.Name,
		"email": req.Email,
	})
}

// newRouter はこのレッスンで使うルーティングを登録したルーターを返します。
func newRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/users", createUserHandler)

	return r
}

func main() {
	newRouter().Run(":8080")
}
