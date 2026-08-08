# Lesson 5（道場）: GORMでCRUDを完成させよう

## 今回学ぶこと

このレッスンは「道場」、つまりCourse「データベース連携」のまとめ問題です。Lesson 4で学んだGORMの基本に加えて、更新・削除・エラーハンドリングまで含めた完全なCRUDを実装します。

## 解説

Lesson 4では `Create` と `Find` を扱いました。ここでは残りの操作を学びます。

### First: 主キーで1件取得する

```go
func getBook(db *gorm.DB, id uint) (*Book, error) {
	var book Book
	if err := db.First(&book, id).Error; err != nil {
		return nil, err // 見つからなければ gorm.ErrRecordNotFound
	}
	return &book, nil
}
```

呼び出し側は `errors.Is(err, gorm.ErrRecordNotFound)` で「見つからなかった」場合を判別できます。

### Save / Updates: 更新する

「取得 → フィールドを書き換え → 保存」という流れが素直で読みやすい書き方です。

```go
func updateBook(db *gorm.DB, id uint, title, author string) (*Book, error) {
	var book Book
	if err := db.First(&book, id).Error; err != nil {
		return nil, err
	}

	book.Title = title
	book.Author = author

	if err := db.Save(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}
```

`Save` は主キーが設定されていれば UPDATE、無ければ INSERT を実行します。

先に `First` で存在確認しているので、存在しないIDに対しては `gorm.ErrRecordNotFound` が返ります。Lesson 3で `RowsAffected` を使って自分でやっていた存在チェックが、この書き方だと自然に実現できます。

### Delete: 削除する

```go
func deleteBook(db *gorm.DB, id uint) error {
	result := db.Delete(&Book{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
```

`Delete` は存在しないIDを指定してもエラーになりません（Lesson 3のSQLと同じ挙動です）。GORMでも `result.RowsAffected` で実際に削除された件数を確認し、0件なら「見つからなかった」として扱います。

> **補足**: GORMには論理削除（`DeletedAt` フィールドを持つモデルでは、実際には削除せず削除日時を記録する）という機能もあります。今回の `Book` は `DeletedAt` を持たないので、物理削除になります。

## 演習

次の4つの関数を実装してください。

**`createBook(db *gorm.DB, title, author string) (*Book, error)`**
- 1件作成し、採番されたIDが入った `*Book` を返す

**`getBook(db *gorm.DB, id uint) (*Book, error)`**
- 主キーで1件取得する。見つからなければGORMのエラー（`gorm.ErrRecordNotFound`）をそのまま返す

**`updateBook(db *gorm.DB, id uint, title, author string) (*Book, error)`**
- 取得 → 書き換え → 保存 の流れで更新し、更新後の `*Book` を返す
- 見つからなければ `gorm.ErrRecordNotFound` を返す

**`deleteBook(db *gorm.DB, id uint) error`**
- 削除する。該当レコードが無ければ `gorm.ErrRecordNotFound` を返す

`Book` 型と `setupDB` はすでに実装済みです。

## ヒント

- `getBook` / `updateBook` は `db.First(&book, id).Error` のエラーをそのまま返せば、自然と `gorm.ErrRecordNotFound` が伝わります
- `deleteBook` だけは `result.RowsAffected == 0` の確認が必要です
