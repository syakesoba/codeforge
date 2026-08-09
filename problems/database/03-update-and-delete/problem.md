# Lesson 3: データを更新・削除する（UPDATE / DELETE）

## 今回学ぶこと

- `UPDATE` / `DELETE` を `db.Exec` で実行する方法
- `RowsAffected()` で「実際に何行変わったか」を確認する方法
- 「対象が存在しなかった」ことをエラーとして扱う書き方

### UPDATE / DELETE も db.Exec

Lesson 2の `INSERT` と同じく、`UPDATE` と `DELETE` も結果の行を返さないので `db.Exec` で実行します。

```go
result, err := db.Exec(`UPDATE books SET title = ?, author = ? WHERE id = ?`, title, author, id)
```

### RowsAffected で存在チェックする

ここが重要なポイントです。存在しないIDに対して `UPDATE` や `DELETE` を実行しても、**SQLとしてはエラーになりません**（「0行が更新された」という成功扱いになります）。

「指定したIDの本が無かった」ことを検知したい場合は、`sql.Result` の `RowsAffected()` で実際に影響を受けた行数を確認します。

```go
var ErrNotFound = errors.New("book not found")

func updateBook(db *sql.DB, id int64, title, author string) error {
	result, err := db.Exec(`UPDATE books SET title = ?, author = ? WHERE id = ?`, title, author, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound // 該当する行が無かった
	}

	return nil
}
```

呼び出し側は `errors.Is(err, ErrNotFound)` で「見つからなかった」場合だけを判別でき、Web APIなら404を返す、といった分岐ができます。Course「Web APIの基礎」やGinコースの道場で404を返していたのと同じ考え方を、DB層で表現しているわけです。

`DELETE` もまったく同じ形です。

```go
func deleteBook(db *sql.DB, id int64) error {
	result, err := db.Exec(`DELETE FROM books WHERE id = ?`, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}

	return nil
}
```

## 演習

`updateBook` と `deleteBook` の2つの関数を実装してください。

**`updateBook(db *sql.DB, id int64, title, author string) error`**
- 指定IDの本の `title` と `author` を更新する
- 該当する行が無ければ `ErrNotFound` を返す

**`deleteBook(db *sql.DB, id int64) error`**
- 指定IDの本を削除する
- 該当する行が無ければ `ErrNotFound` を返す

`ErrNotFound` と `setupDB` / `insertBook` はすでに用意されています。

## ヒント

- どちらも「Exec → RowsAffected → 0なら ErrNotFound → それ以外は nil」という同じ流れです
- `RowsAffected()` の戻り値は `(int64, error)` の2つです
