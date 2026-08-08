# Lesson 2: ルーティンググループとパスパラメータ

## 今回学ぶこと

- `r.Group` で共通のパス接頭辞（プレフィックス）をまとめる方法
- `c.Param` でパスパラメータを取得する方法
- `c.Query` でクエリパラメータを取得する方法

## 解説

### ルーティンググループ

実際のAPIでは `/api/v1/users`, `/api/v1/items` のように、共通の接頭辞を持つエンドポイントがたくさん並びます。毎回フルパスを書くのは冗長なので、Ginでは `r.Group` でまとめられます。

```go
func newRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", listUsersHandler)   // 実際のパスは /api/v1/users
		v1.GET("/items", listItemsHandler)   // 実際のパスは /api/v1/items
	}

	return r
}
```

`{ }` のブロックは文法上必須ではありませんが、「このグループに属するルート」を視覚的にまとめるためにGinのドキュメントでもよく使われる書き方です。

グループは入れ子にもできます。Lesson 4で学ぶミドルウェアを「このグループにだけ適用する」といった使い方もでき、認証が必要なエンドポイント群をまとめるときに便利です。

### パスパラメータ

Ginでは、パスの中で `:名前` と書くとパスパラメータになります。ハンドラの中では `c.Param("名前")` で取り出します。

```go
v1.GET("/users/:id", func(c *gin.Context) {
	id := c.Param("id") // /api/v1/users/42 なら "42"
	c.JSON(http.StatusOK, gin.H{"id": id})
})
```

Course「Web APIの基礎」の `net/http` では `{id}` と書いて `r.PathValue("id")` で取り出しましたが、Ginでは `:id` と `c.Param("id")` になります。書き方が違うだけで考え方は同じです。

### クエリパラメータ

`/search?q=golang` のようなクエリパラメータは `c.Query` で取得します。

```go
q := c.Query("q")                    // 無ければ空文字列 ""
limit := c.DefaultQuery("limit", "10") // 無ければ "10"
```

## 演習

`getUserHandler` 関数を実装してください。パスパラメータ `id` を受け取り、次のJSONを**ステータスコード200**で返します。

```json
{"id": "42", "name": "user-42"}
```

（`/api/v1/users/42` にリクエストが来た場合の例です。`name` は `"user-"` にidを繋げた文字列にしてください）

`newRouter()` の中で `GET /api/v1/users/:id` はすでに `getUserHandler` に紐づけて登録されているので、あなたが実装するのはハンドラの中身だけです。

## ヒント

- `id := c.Param("id")` でパスパラメータを取得できます
- 文字列の連結は `"user-" + id` で書けます
- `c.JSON(http.StatusOK, gin.H{"id": id, "name": "user-" + id})` の形になります
