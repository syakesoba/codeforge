# Lesson 1: Ginの基本

## 今回学ぶこと

- Webフレームワーク **Gin** の基本的な使い方
- `gin.Engine` にハンドラを登録する方法
- `c.JSON` でJSONレスポンスを返す方法

### Ginとは

Course「Web APIの基礎」では標準ライブラリの `net/http` だけでAPIを作りました。Ginは、それをもっと短く・読みやすく書けるようにするWebフレームワークです。Goのフレームワークの中で最も広く使われており、実務でもよく登場します。

`net/http` で書いていたハンドラは、Ginでは次のように書けます。

```go
// net/http の場合
func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
}

// Gin の場合
func pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
```

Ginでは `http.ResponseWriter` と `*http.Request` の代わりに、両方をまとめた `*gin.Context` を1つ受け取ります。`c.JSON(ステータスコード, 値)` を呼ぶだけで、`Content-Type` の設定・ステータスコードの書き込み・JSONエンコードがすべて行われます。

`gin.H` は `map[string]any` のエイリアスで、その場限りのJSONを組み立てるときに便利です。もちろん構造体を渡すこともできます。

### ルーターを組み立てる

Ginでは `gin.Engine` がルーター（`net/http` の `*http.ServeMux` に相当）です。

```go
func newRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.GET("/ping", pingHandler)

	return r
}

func main() {
	newRouter().Run(":8080")
}
```

- `gin.New()`: 何もついていない素のルーターを作る（`gin.Default()` はログとリカバリのミドルウェア付き）
- `r.GET(パス, ハンドラ)`: GETリクエストのルートを登録する。`POST` / `PUT` / `DELETE` も同様
- `r.Run(":8080")`: サーバーを起動する

> **補足**: このコースの `newRouter()` では最初に `gin.SetMode(gin.TestMode)` を呼んでいます。これを書かないとGinが大量のデバッグログを出力するためです。実務のコードでは通常書きません（環境変数 `GIN_MODE=release` などで制御します）。

`*gin.Engine` は `http.Handler` を満たしているので、Course「Web APIの基礎」で使った `httptest` でそのままテストできます。

## 演習

`healthHandler` 関数を実装してください。次のJSONを**ステータスコード200**で返します。

```json
{"status": "ok"}
```

`newRouter()` の中で `GET /health` はすでに `healthHandler` に紐づけて登録されているので、あなたが実装するのはハンドラの中身だけです。

## ヒント

- `c.JSON(http.StatusOK, gin.H{"status": "ok"})` の1行で書けます
