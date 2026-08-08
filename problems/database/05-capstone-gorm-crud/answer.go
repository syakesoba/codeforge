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

// getBook は主キーで1件取得します。見つからなければ gorm.ErrRecordNotFound を返します。
func getBook(db *gorm.DB, id uint) (*Book, error) {
	var book Book
	if err := db.First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

// updateBook は指定IDの本を更新し、更新後の *Book を返します。
// 見つからなければ gorm.ErrRecordNotFound を返します。
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

// deleteBook は指定IDの本を削除します。
// 該当レコードが無ければ gorm.ErrRecordNotFound を返します。
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

func main() {}
