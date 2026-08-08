package main

import (
	"github.com/gin-gonic/gin"
)

// createUserRequest はリクエストボディ {"name": "...", "email": "..."} を受け取る型です。
type createUserRequest struct {
	// TODO: binding タグを追加してください
	//   Name  は必須
	//   Email は必須かつメールアドレス形式
	Name  string `json:"name"`
	Email string `json:"email"`
}

// createUserHandler はリクエストボディをバインド・バリデーションし、
// 失敗なら400、成功なら201でレスポンスを返します。
func createUserHandler(c *gin.Context) {
	// TODO:
	// 1. var req createUserRequest を用意し、c.ShouldBindJSON(&req) でバインドする
	// 2. エラーなら c.JSON(400, gin.H{"error": err.Error()}) を返して return する
	// 3. 成功なら c.JSON(201, gin.H{"name": ..., "email": ...}) を返す
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
