package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// このテストは、レッスンの採点で実際に使われるパッケージ
// （testing / httptest / gin のルーター実行経路）のビルドキャッシュを
// イメージに焼き込むためだけに存在する。
// これが無いと、Ginを使うレッスンの「初回提出」だけテストバイナリの
// フルコンパイルが走り、タイムアウトしてしまう。
func TestPriming(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("priming request failed: %d", rec.Code)
	}
}
