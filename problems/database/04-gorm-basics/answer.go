//go:build ignore

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
	book := Book{Title: title, Author: author}
	if err := db.Create(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

// listBooks は全件を取得して返します。
func listBooks(db *gorm.DB) ([]Book, error) {
	var books []Book
	if err := db.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func main() {}
