# Lesson 4: Ginのミドルウェア

## 今回学ぶこと

- Ginにおけるミドルウェアの書き方（`gin.HandlerFunc`）
- `c.Next()` の役割
- `c.Abort()` / `c.AbortWithStatusJSON` でリクエストを中断する方法
- `r.Use` と グループ単位のミドルウェア適用

### Ginのミドルウェア

Course「Web APIの基礎」では、ミドルウェアを「`http.Handler` を受け取って `http.Handler` を返す関数」として書きました。Ginではもっとシンプルで、**ミドルウェアもハンドラと同じ `func(c *gin.Context)`** です。

```go
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Request-ID", "abc-123")
		c.Next() // 後続のハンドラを実行する
	}
}
```

- 戻り値の型 `gin.HandlerFunc` は `func(c *gin.Context)` のエイリアスです
- `c.Next()` を呼ぶと、後続のミドルウェアや本来のハンドラが実行されます
- `c.Next()` の**前**に書いた処理はハンドラの前に、**後**に書いた処理はハンドラの後に実行されます

### ミドルウェアを適用する

ルーター全体に適用する場合は `r.Use`、特定のグループにだけ適用する場合はグループに対して `Use` を呼びます。

```go
r := gin.New()
r.Use(RequestIDMiddleware()) // 全ルートに適用

admin := r.Group("/admin")
admin.Use(AuthMiddleware())  // /admin 配下だけに適用
```

### リクエストを中断する

認証ミドルウェアのように「条件を満たさなければ後続の処理をさせたくない」場合は、`c.Next()` を呼ばずに `c.AbortWithStatusJSON` で打ち切ります。

```go
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "Bearer secret-token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}
```

`c.AbortWithStatusJSON(コード, 値)` は「そのステータスコードとJSONを返し、後続のハンドラを実行しない」という意味です。必ず `return` とセットで書いてください。

> **補足**: `c.Abort()` を呼ばずに単に `return` しただけでは、Ginは後続のハンドラを実行してしまいます。中断したいときは必ず `Abort` 系のメソッドを使ってください。

## 演習

`APIKeyMiddleware` 関数を実装してください。次の動作をするミドルウェアです。

1. リクエストヘッダー `X-API-Key` の値を取得する（`c.GetHeader("X-API-Key")`）
2. 値が `"codeforge-key"` と**一致しない**場合は、ステータスコード **401** で `{"error": "invalid api key"}` を返し、後続の処理を中断する
3. 一致する場合は `c.Next()` を呼んで後続の処理を続行する

`newRouter()` の中で、このミドルウェアは `/secret` エンドポイントに適用済みです。

## ヒント

- 関数のシグネチャは `func APIKeyMiddleware() gin.HandlerFunc` です。`return func(c *gin.Context) { ... }` の形で中身を書きます
- 中断には `c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})` を使い、その直後に `return` してください
