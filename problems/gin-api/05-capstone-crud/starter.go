package main

import (
	"sync"

	"github.com/gin-gonic/gin"
)

// Book は1冊の書籍を表します。
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var (
	mu     sync.Mutex
	books  = map[int]Book{}
	nextID = 1
)

// listBooksHandler は全書籍をJSON配列で返します（200）。
func listBooksHandler(c *gin.Context) {
	// TODO: books の中身をスライスに詰め直して c.JSON(200, result) で返す
}

// createBookHandler は新しい書籍を作成します（成功: 201 / バリデーションエラー: 400）。
func createBookHandler(c *gin.Context) {
	// TODO:
	// 1. createBookRequest 型を定義し、title/author に binding:"required" を付ける
	// 2. c.ShouldBindJSON でバインドし、失敗したら400を返す
	// 3. nextID を採番して books に登録し、201で作成したBookを返す
}

// updateBookHandler は指定IDの書籍を更新します（成功: 200 / 見つからない: 404）。
func updateBookHandler(c *gin.Context) {
	// TODO:
	// 1. c.Param("id") を strconv.Atoi で数値に変換する（失敗したら400）
	// 2. リクエストボディをバインドする（失敗したら400）
	// 3. 該当する書籍が無ければ404、あれば内容を上書きして200で返す
}

// deleteBookHandler は指定IDの書籍を削除します（成功: 204 / 見つからない: 404）。
func deleteBookHandler(c *gin.Context) {
	// TODO:
	// 1. c.Param("id") を strconv.Atoi で数値に変換する（失敗したら400）
	// 2. 該当する書籍が無ければ404、あれば delete して c.Status(204) を返す
}

// newRouter はこのレッスンで使うルーティングを登録したルーターを返します。
func newRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	api := r.Group("/api")
	{
		api.GET("/books", listBooksHandler)
		api.POST("/books", createBookHandler)
		api.PUT("/books/:id", updateBookHandler)
		api.DELETE("/books/:id", deleteBookHandler)
	}

	return r
}

func main() {
	newRouter().Run(":8080")
}
