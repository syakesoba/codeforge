# Lesson 2: データを登録・取得する（INSERT / SELECT）

## 今回学ぶこと

- `db.Exec` でデータを登録し、採番されたIDを取得する方法
- `db.QueryRow` + `Scan` で1件だけ取得する方法
- `db.Query` + `rows.Next()` で複数件取得する方法
- プレースホルダ（`?`）を使う理由

## 解説

### プレースホルダを必ず使う

SQLに値を埋め込むときは、文字列連結ではなく**プレースホルダ `?`** を使い、値は引数として渡します。

```go
// 良い例
db.Exec(`INSERT INTO books (title, author) VALUES (?, ?)`, title, author)

// 悪い例（SQLインジェクションの脆弱性になる）
db.Exec("INSERT INTO books (title, author) VALUES ('" + title + "', '" + author + "')")
```

文字列連結だと、ユーザーが入力した値にSQL文の一部を紛れ込ませることで、データを盗まれたり消されたりする「SQLインジェクション」という攻撃を受けます。プレースホルダを使えばドライバが安全にエスケープしてくれるので、**必ずこちらを使ってください**。

### 登録して採番されたIDを取得する

`db.Exec` の戻り値 `sql.Result` から `LastInsertId()` で、自動採番されたIDを取得できます。

```go
func insertBook(db *sql.DB, title, author string) (int64, error) {
	result, err := db.Exec(`INSERT INTO books (title, author) VALUES (?, ?)`, title, author)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
```

### 1件だけ取得する（QueryRow + Scan）

結果が1行だけとわかっている場合は `db.QueryRow` を使い、`Scan` で変数に流し込みます。

```go
func getBookTitle(db *sql.DB, id int64) (string, error) {
	var title string
	err := db.QueryRow(`SELECT title FROM books WHERE id = ?`, id).Scan(&title)
	if err != nil {
		return "", err
	}
	return title, nil
}
```

`Scan` には**ポインタ**を渡します（`&title`）。SELECTしたカラムの順番と、`Scan` に渡す変数の順番・型を一致させる必要があります。

該当する行が無い場合、`Scan` は `sql.ErrNoRows` というエラーを返します。「見つからなかった」と「本当のエラー」を区別したいときは `errors.Is(err, sql.ErrNoRows)` で判定します。

### 複数件取得する（Query + rows.Next）

複数行を扱う場合は `db.Query` を使い、`rows.Next()` でループします。

```go
func listBooks(db *sql.DB) ([]Book, error) {
	rows, err := db.Query(`SELECT id, title, author FROM books ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // 必ず閉じる

	var books []Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author); err != nil {
			return nil, err
		}
		books = append(books, b)
	}

	// ループ中に発生したエラーはここで確認する
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}
```

押さえるべきポイントは3つです。

1. `defer rows.Close()` を必ず書く（書かないとDB接続が解放されない）
2. `for rows.Next()` の中で `rows.Scan` して1行ずつ取り出す
3. ループを抜けた後に `rows.Err()` を確認する（ループ中に起きたエラーはここに入る）

## 演習

`insertBook` と `listBooks` の2つの関数を実装してください。

**`insertBook(db *sql.DB, title, author string) (int64, error)`**
- `books` テーブルに1件登録し、採番されたIDを返す

**`listBooks(db *sql.DB) ([]Book, error)`**
- `books` テーブルの全件を **id順** で取得して返す

`Book` 型とテーブルを用意する `setupDB` はすでに実装済みです。

## ヒント

- `insertBook` は `result.LastInsertId()` の戻り値をそのまま返せます
- `listBooks` の SELECT は `SELECT id, title, author FROM books ORDER BY id` です
- `rows.Scan(&b.ID, &b.Title, &b.Author)` のように、SELECTした順番どおりにポインタを渡します
