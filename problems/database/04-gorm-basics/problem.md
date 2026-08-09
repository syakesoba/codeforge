# Lesson 4: GORMの基本

## 今回学ぶこと

- ORM（Object-Relational Mapping）とは何か
- GORMでの接続と `AutoMigrate`
- `Create` / `First` / `Find` による基本操作

### ORMとは

Lesson 1〜3では、SQL文を自分で書き、`rows.Scan` で1つずつ変数に詰めていました。動作は明快ですが、テーブルが増えると同じような定型コードが大量に発生します。

**ORM**は、Goの構造体とDBのテーブルを対応づけて、SQLを書かずにDB操作できるようにするライブラリです。Goで最も広く使われているORMが **GORM** です。

### モデルを定義する

GORMでは、構造体がそのままテーブルの定義になります。

```go
type Book struct {
	ID     uint   `gorm:"primaryKey"`
	Title  string `gorm:"not null"`
	Author string `gorm:"not null"`
}
```

- `gorm:"primaryKey"` で主キーを指定します（`ID` という名前のフィールドは指定しなくても自動で主キーになります）
- テーブル名は構造体名の複数形（`Book` → `books`）が自動的に使われます

### 接続とAutoMigrate

```go
import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Book{}); err != nil {
		return nil, err
	}

	return db, nil
}
```

`AutoMigrate` は、構造体の定義に合わせてテーブルを自動的に作成・更新してくれる機能です。Lesson 1で自分で書いた `CREATE TABLE` が不要になります。

> **注意**: Lesson 1〜3で使ったドライバは `github.com/glebarez/go-sqlite`（`database/sql` 用）でしたが、GORMで使うのは `github.com/glebarez/sqlite`（GORM用）です。名前が似ているので混同しないよう注意してください。

### 基本操作

```go
// 作成: 構造体のポインタを渡す。採番されたIDは book.ID に書き戻される
book := Book{Title: "Go入門", Author: "Gopher"}
if err := db.Create(&book).Error; err != nil {
	return err
}
fmt.Println(book.ID) // 1

// 主キーで1件取得
var found Book
if err := db.First(&found, id).Error; err != nil {
	return err // 見つからない場合は gorm.ErrRecordNotFound
}

// 全件取得
var books []Book
if err := db.Find(&books).Error; err != nil {
	return err
}
```

GORMのメソッドは `*gorm.DB` を返し、その `.Error` フィールドにエラーが入ります。`if err := db.Create(&book).Error; err != nil` という書き方に慣れてください。

`First` で該当レコードが無い場合は `gorm.ErrRecordNotFound` というエラーが返るので、`errors.Is(err, gorm.ErrRecordNotFound)` で判定できます。Lesson 3で自作した `ErrNotFound` を、GORMが最初から用意してくれているわけです。

## 演習

`createBook` と `listBooks` の2つの関数を実装してください。

**`createBook(db *gorm.DB, title, author string) (*Book, error)`**
- `Book` を1件作成し、採番されたIDが入った状態の `*Book` を返す

**`listBooks(db *gorm.DB) ([]Book, error)`**
- 全件を取得して返す

`Book` 型と `setupDB` はすでに実装済みです。

## ヒント

- `createBook` は `book := Book{Title: title, Author: author}` を作って `db.Create(&book)` し、`&book, nil` を返します
- `listBooks` は `var books []Book` を宣言して `db.Find(&books)` するだけです
- どちらも `.Error` を確認するのを忘れないでください
