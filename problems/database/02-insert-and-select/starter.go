package main

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

// Book は1冊の書籍を表します。
type Book struct {
	ID     int64
	Title  string
	Author string
}

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

// insertBook は books テーブルに1件登録し、採番されたIDを返します。
func insertBook(db *sql.DB, title, author string) (int64, error) {
	// TODO:
	// 1. db.Exec でプレースホルダ(?)を使ってINSERTする
	// 2. 戻り値の sql.Result から LastInsertId() を返す
	return 0, nil
}

// listBooks は books テーブルの全件を id順 で取得して返します。
func listBooks(db *sql.DB) ([]Book, error) {
	// TODO:
	// 1. db.Query で SELECT id, title, author FROM books ORDER BY id を実行する
	// 2. defer rows.Close() を忘れずに書く
	// 3. for rows.Next() で1行ずつ rows.Scan して []Book に詰める
	// 4. ループ後に rows.Err() を確認してから返す
	return nil, nil
}

func main() {}
