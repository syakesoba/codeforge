package main

import (
	"database/sql"
	"errors"

	_ "github.com/glebarez/go-sqlite"
)

// ErrNotFound は指定したIDの本が存在しなかったことを表すエラーです。
var ErrNotFound = errors.New("book not found")

// setupDB はインメモリSQLiteに接続し、books テーブルを作成して返します（実装済み）。
func setupDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE books (
			id     INTEGER PRIMARY KEY AUTOINCREMENT,
			title  TEXT NOT NULL,
			author TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// insertBook は books テーブルに1件登録し、採番されたIDを返します（実装済み）。
func insertBook(db *sql.DB, title, author string) (int64, error) {
	result, err := db.Exec(`INSERT INTO books (title, author) VALUES (?, ?)`, title, author)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// updateBook は指定IDの本の title と author を更新します。
// 該当する行が無ければ ErrNotFound を返します。
func updateBook(db *sql.DB, id int64, title, author string) error {
	// TODO:
	// 1. db.Exec で UPDATE books SET title = ?, author = ? WHERE id = ? を実行する
	// 2. result.RowsAffected() で影響を受けた行数を確認する
	// 3. 0行なら ErrNotFound を返し、それ以外は nil を返す
	return nil
}

// deleteBook は指定IDの本を削除します。
// 該当する行が無ければ ErrNotFound を返します。
func deleteBook(db *sql.DB, id int64) error {
	// TODO:
	// 1. db.Exec で DELETE FROM books WHERE id = ? を実行する
	// 2. result.RowsAffected() で影響を受けた行数を確認する
	// 3. 0行なら ErrNotFound を返し、それ以外は nil を返す
	return nil
}

func main() {}
