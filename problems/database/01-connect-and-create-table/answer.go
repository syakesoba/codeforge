//go:build ignore

package main

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

// setupDB はインメモリSQLiteに接続し、books テーブルを作成して返します。
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

func main() {}
