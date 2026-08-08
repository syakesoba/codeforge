//go:build ignore

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
	result, err := db.Exec(`INSERT INTO books (title, author) VALUES (?, ?)`, title, author)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// listBooks は books テーブルの全件を id順 で取得して返します。
func listBooks(db *sql.DB) ([]Book, error) {
	rows, err := db.Query(`SELECT id, title, author FROM books ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author); err != nil {
			return nil, err
		}
		books = append(books, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

func main() {}
