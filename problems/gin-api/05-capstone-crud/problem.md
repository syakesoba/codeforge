# Lesson 5（道場）: GinでCRUD APIを作ろう

## 今回学ぶこと

このレッスンは「道場」、つまりCourse「フレームワークで作るAPI」のまとめ問題です。Lesson 1〜4で学んだGinの機能をすべて使って、書籍を管理するCRUD APIを作ります。

Course「Web APIの基礎」の道場では同じようなAPIを `net/http` で作りました。今回はそれをGinで書き直すことで、フレームワークを使うとどれだけ簡潔になるかを体感できます。

## 解説: 作るものの仕様

| メソッド | パス | 役割 | 成功時のステータス |
|---|---|---|---|
| `GET` | `/api/books` | 全書籍を一覧で返す | 200 |
| `POST` | `/api/books` | 新しい書籍を作成する | 201 |
| `PUT` | `/api/books/:id` | 指定IDの書籍を更新する | 200 |
| `DELETE` | `/api/books/:id` | 指定IDの書籍を削除する | 204 |

データ構造と保存先はすでに用意されています。

```go
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
```

`newRouter()` の中で、4つのエンドポイントは `/api` グループ配下にすでに登録済みです。あなたが実装するのは4つのハンドラの中身だけです。

### listBooksHandler（一覧取得）

`books` の中身をスライスに詰め直して200で返します。

```go
func listBooksHandler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	result := make([]Book, 0, len(books))
	for _, b := range books {
		result = append(result, b)
	}

	c.JSON(http.StatusOK, result)
}
```

### createBookHandler（作成）

`binding` タグでバリデーションし（Lesson 3）、失敗なら400、成功なら201を返します。

```go
type createBookRequest struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

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
```

### updateBookHandler（更新）

パスパラメータ `:id` を `strconv.Atoi` で数値に変換し（Lesson 2）、該当する書籍があれば更新して200、なければ404を返します。

```go
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
```

### deleteBookHandler（削除）

該当する書籍があれば削除して**204（本文なし）**、なければ404を返します。

```go
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
```

204は本文を返さないステータスコードなので、`c.JSON` ではなく `c.Status` を使う点に注意してください。

## 演習

上記の仕様どおりに、4つのハンドラ関数（`listBooksHandler`, `createBookHandler`, `updateBookHandler`, `deleteBookHandler`）と、リクエスト用の構造体 `createBookRequest` を実装してください。

採点では、作成 → 一覧確認 → 更新 → 存在しないIDへの更新（404）→ バリデーションエラー（400）→ 削除 → 一覧確認 → 存在しないIDへの削除（404）という一連の流れを自動でテストします。

## ヒント

- `createBookRequest` は作成・更新の両方で使い回せます
- `mu.Lock()` / `defer mu.Unlock()` で共有マップを保護してください（Course「ゴルーチン・並行処理」Lesson 3を参照）
