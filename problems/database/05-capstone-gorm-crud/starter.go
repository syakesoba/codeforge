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
	// TODO: db.Create(&book) で作成し、&book を返す
	return nil, nil
}

// getBook は主キーで1件取得します。見つからなければ gorm.ErrRecordNotFound を返します。
func getBook(db *gorm.DB, id uint) (*Book, error) {
	// TODO: db.First(&book, id) で取得し、エラーはそのまま返す
	return nil, nil
}

// updateBook は指定IDの本を更新し、更新後の *Book を返します。
// 見つからなければ gorm.ErrRecordNotFound を返します。
func updateBook(db *gorm.DB, id uint, title, author string) (*Book, error) {
	// TODO:
	// 1. db.First(&book, id) で取得する（エラーならそのまま返す）
	// 2. book.Title / book.Author を書き換える
	// 3. db.Save(&book) で保存し、&book を返す
	return nil, nil
}

// deleteBook は指定IDの本を削除します。
// 該当レコードが無ければ gorm.ErrRecordNotFound を返します。
func deleteBook(db *gorm.DB, id uint) error {
	// TODO:
	// 1. result := db.Delete(&Book{}, id) を実行する
	// 2. result.Error を確認する
	// 3. result.RowsAffected == 0 なら gorm.ErrRecordNotFound を返す
	return nil
}

func main() {}
