//go:build ignore

package main

import (
	"net/http"
	"strconv"
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

type createBookRequest struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

// listBooksHandler は全書籍をJSON配列で返します（200）。
func listBooksHandler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	result := make([]Book, 0, len(books))
	for _, b := range books {
		result = append(result, b)
	}

	c.JSON(http.StatusOK, result)
}

// createBookHandler は新しい書籍を作成します（成功: 201 / バリデーションエラー: 400）。
func createBookHandler(c *gin.Context) {
	var req createBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	book := Book{ID: nextID, Title: req.Title, Author: req.Author}
	books[nextID] = book
	nextID++
	mu.Unlock()

	c.JSON(http.StatusCreated, book)
}

// updateBookHandler は指定IDの書籍を更新します（成功: 200 / 見つからない: 404）。
func updateBookHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req createBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	book, ok := books[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	book.Title = req.Title
	book.Author = req.Author
	books[id] = book

	c.JSON(http.StatusOK, book)
}

// deleteBookHandler は指定IDの書籍を削除します（成功: 204 / 見つからない: 404）。
func deleteBookHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, ok := books[id]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	delete(books, id)

	c.Status(http.StatusNoContent)
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
