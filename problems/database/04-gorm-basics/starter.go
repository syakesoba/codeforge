package main

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// Book は books テーブルに対応するモデルです。
type Book struct {
	ID     uint   `gorm:"primaryKey"`
	Title  string `gorm:"not null"`
	Author string `gorm:"not null"`
}

// setupDB はインメモリSQLiteに接続し、AutoMigrateでテーブルを作成して返します（実装済み）。
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

// createBook は Book を1件作成し、採番されたIDが入った *Book を返します。
func createBook(db *gorm.DB, title, author string) (*Book, error) {
	// TODO:
	// 1. book := Book{Title: title, Author: author} を作る
	// 2. db.Create(&book).Error を確認する（エラーなら nil, err を返す）
	// 3. &book, nil を返す
	return nil, nil
}

// listBooks は全件を取得して返します。
func listBooks(db *gorm.DB) ([]Book, error) {
	// TODO:
	// 1. var books []Book を宣言する
	// 2. db.Find(&books).Error を確認する（エラーなら nil, err を返す）
	// 3. books, nil を返す
	return nil, nil
}

func main() {}
