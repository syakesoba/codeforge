# Lesson 4: 認証ミドルウェアでAPIを保護する

## 今回学ぶこと

- `Authorization: Bearer <token>` ヘッダーの形式
- JWT検証をミドルウェアとして組み込む方法
- `c.Set` / `c.Get` でハンドラに認証情報を引き渡す方法

### Bearerトークン

クライアントはJWTを `Authorization` ヘッダーに次の形式で載せて送ります。

```
Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJhbGljZSJ9.xxxxx
```

`Bearer `（末尾に半角スペース）というプレフィックスが付くのが決まりです。サーバー側ではこのプレフィックスを取り除いてからトークンを検証します。

`strings.CutPrefix` を使うと、「プレフィックスを取り除く」と「そもそもプレフィックスが付いていたか」を同時に判定できて便利です。

```go
tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
if !ok {
	// "Bearer " で始まっていない → 不正な形式
}
```

### 認証ミドルウェア

GinコースのLesson 4で学んだミドルウェアの形に、JWT検証を組み込みます。

```go
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		username, err := VerifyToken(secret, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		c.Set("username", username) // 後続のハンドラに引き渡す
		c.Next()
	}
}
```

ポイントは次の3つです。

1. ヘッダーの形式が不正なら401で中断する
2. トークンの検証に失敗（署名不正・期限切れなど）しても401で中断する
3. 成功したら `c.Set` でユーザー名を保存し、`c.Next()` で後続に進む

> **セキュリティ上の注意**: エラーの理由（「ヘッダーが無い」「期限切れ」「署名が不正」など）を細かく返すと、攻撃者に手がかりを与えることがあります。認証エラーは一律 `unauthorized` として返すのが無難です。

### c.Set と c.Get

`c.Set(key, value)` でリクエストのコンテキストに値を保存でき、後続のハンドラで `c.Get(key)` として取り出せます。

```go
func meHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "username not found in context"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": username})
}
```

`c.Get` の戻り値は `(any, bool)` です。型が確定している場合は `c.GetString("username")` のような専用メソッドも使えます。

これにより「認証はミドルウェアが担当し、ハンドラは認証済みである前提で本来の処理に集中する」という責務の分離ができます。

## 演習

`AuthMiddleware` 関数を実装してください。

**`AuthMiddleware(secret string) gin.HandlerFunc`**
1. `Authorization` ヘッダーを取得する
2. `"Bearer "` プレフィックスが付いていなければ、**401** で `{"error": "unauthorized"}` を返して中断する
3. `VerifyToken` でトークンを検証し、失敗したら同様に **401** で中断する
4. 成功したら `c.Set("username", username)` してから `c.Next()` を呼ぶ

`IssueToken` / `VerifyToken` と、保護対象の `meHandler` はすでに実装済みです。

## ヒント

- `strings.CutPrefix(authHeader, "Bearer ")` の戻り値は `(string, bool)` です
- 中断には `c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})` を使い、直後に `return` してください
