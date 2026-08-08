# Lesson 5（道場）: ログインAPIを作ろう

## 今回学ぶこと

このレッスンは「道場」、つまりCourse「認証・JWT」のまとめ問題です。Lesson 1〜4で学んだすべて（bcrypt・JWT発行・JWT検証・認証ミドルウェア）を組み合わせて、実際に動くログインAPIを作ります。

## 解説: 作るものの仕様

| メソッド | パス | 認証 | 役割 |
|---|---|---|---|
| `POST` | `/register` | 不要 | ユーザーを登録する |
| `POST` | `/login` | 不要 | 認証してJWTを発行する |
| `GET` | `/profile` | **必要** | 認証済みユーザーの情報を返す |

ユーザーの保存先とヘルパー関数はすでに用意されています。

```go
type User struct {
	Username     string
	PasswordHash string
}

var (
	mu    sync.Mutex
	users = map[string]User{} // key: username
)
```

`HashPassword` / `CheckPassword` / `IssueToken` / `VerifyToken` / `AuthMiddleware` はすべて実装済みです。あなたが実装するのは3つのハンドラだけです。

### registerHandler（ユーザー登録）

リクエストボディ `{"username": "...", "password": "..."}` を受け取ります。

```go
type credentialsRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

func registerHandler(c *gin.Context) {
	var req credentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := users[req.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
		return
	}

	users[req.Username] = User{Username: req.Username, PasswordHash: hash}

	c.JSON(http.StatusCreated, gin.H{"username": req.Username})
}
```

ポイントは3つです。

- パスワードは `min=8` で最低8文字を要求する（バリデーションはGinコースLesson 3の復習）
- **平文パスワードは絶対に保存しない**。必ず `HashPassword` を通す
- 既に同じユーザー名が存在する場合は **409 Conflict** を返す

### loginHandler（ログイン）

```go
func loginHandler(c *gin.Context) {
	var req credentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	user, exists := users[req.Username]
	mu.Unlock()

	if !exists || !CheckPassword(user.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := IssueToken(jwtSecret, user.Username, time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
```

セキュリティ上の重要なポイントとして、**「ユーザーが存在しない」と「パスワードが違う」を区別せず**、どちらも同じ `invalid credentials` として401を返しています。区別してしまうと、攻撃者に「このユーザー名は存在する」という情報を与えてしまうためです。

### profileHandler（認証必須）

`AuthMiddleware` が `c.Set("username", ...)` で保存した値を取り出すだけです。

```go
func profileHandler(c *gin.Context) {
	username := c.GetString("username")
	c.JSON(http.StatusOK, gin.H{"username": username})
}
```

`c.GetString` は `c.Get` の文字列版で、型アサーションを省けます。

## 演習

上記の仕様どおりに、3つのハンドラ（`registerHandler`, `loginHandler`, `profileHandler`）と、リクエスト用の構造体 `credentialsRequest` を実装してください。

採点では、登録 → 重複登録（409）→ 短すぎるパスワード（400）→ ログイン成功 → 誤ったパスワード（401）→ 存在しないユーザー（401）→ トークンで `/profile` にアクセス → 認証なしで `/profile`（401）という一連の流れをテストします。

## ヒント

- `credentialsRequest` は登録・ログインの両方で使い回せます
- JWT発行に使う秘密鍵はパッケージ変数 `jwtSecret` に用意してあります
- 有効期限は `time.Hour` を使ってください
