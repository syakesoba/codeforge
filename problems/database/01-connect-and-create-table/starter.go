package main

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

// setupDB はインメモリSQLiteに接続し、books テーブルを作成して返します。
func setupDB() (*sql.DB, error) {
	// TODO:
	// 1. sql.Open("sqlite", ":memory:") でDBを開く（エラーなら nil, err を返す）
	// 2. db.Exec で books テーブルを作成する
	//    id: INTEGER PRIMARY KEY AUTOINCREMENT
	//    title: TEXT NOT NULL
	//    author: TEXT NOT NULL
	// 3. 成功したら db, nil を返す
	return nil, nil
}

func main() {}
