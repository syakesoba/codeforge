# Lesson 3: リクエストのバインディングとバリデーション

## 今回学ぶこと

- `c.ShouldBindJSON` でリクエストボディを構造体に流し込む方法
- `binding` タグによる入力バリデーション
- バリデーションエラーを 400 Bad Request として返す方法

## 解説

### ShouldBindJSON

Course「Web APIの基礎」では `json.NewDecoder(r.Body).Decode(&req)` でリクエストボディを読み取りました。Ginでは `c.ShouldBindJSON(&req)` の1行で同じことができます。

```go
type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func createUserHandler(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"name": req.Name})
}
```

### bindingタグでバリデーションする

Ginの真価はここからです。構造体のフィールドに `binding` タグを付けると、`ShouldBindJSON` が**バリデーションまで自動で行ってくれます**。

```go
type createUserRequest struct {
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age"   binding:"gte=0,lte=150"`
}
```

よく使うルールは次のとおりです。

| タグ | 意味 |
|---|---|
| `required` | 必須。ゼロ値（空文字列や0）だとエラー |
| `email` | メールアドレスの形式であること |
| `min=3` / `max=20` | 文字列なら長さ、数値なら値の下限・上限 |
| `gte=0` / `lte=150` | 数値の下限・上限（以上・以下） |

複数のルールはカンマで繋げます（`binding:"required,email"`）。

バリデーションに失敗すると `ShouldBindJSON` がエラーを返すので、それを `400 Bad Request` として返すのが定石です。自分で `if req.Name == "" { ... }` と書かなくてよいのがGinを使う大きなメリットです。

### 実装のかたち

```go
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
```

「バインド → 失敗したら400を返して `return` → 成功したら本来の処理」という流れを覚えてください。`return` を忘れるとエラーレスポンスの後に正常系の処理も走ってしまうので注意が必要です。

## 演習

`createUserHandler` 関数と、リクエスト用の構造体 `createUserRequest` を完成させてください。

**構造体のバリデーション要件**（`binding` タグを追加してください）:
- `Name`: 必須
- `Email`: 必須かつメールアドレス形式

**ハンドラの動作**:
- バインド・バリデーションに失敗したら、ステータスコード **400** で `{"error": "<エラーメッセージ>"}` を返す
- 成功したら、ステータスコード **201** で `{"name": "<name>", "email": "<email>"}` を返す

## ヒント

- 構造体タグは `` `json:"name" binding:"required"` `` のようにスペース区切りで並べます
- エラーメッセージは `err.Error()` をそのまま使って構いません
