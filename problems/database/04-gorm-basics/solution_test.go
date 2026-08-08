package main

import (
	"testing"

	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := setupDB()
	if err != nil {
		t.Fatalf("setupDB がエラーを返しました: %v", err)
	}
	return db
}

func TestCreateBookReturnsBookWithID(t *testing.T) {
	db := newTestDB(t)

	book, err := createBook(db, "Go入門", "Gopher")
	if err != nil {
		t.Fatalf("createBook がエラーを返しました: %v", err)
	}
	if book == nil {
		t.Fatal("createBook が nil を返しました。&book を返していますか？")
	}
	if book.ID == 0 {
		t.Fatal("IDが採番されていません。db.Create(&book) のようにポインタを渡していますか？")
	}
	if book.Title != "Go入門" || book.Author != "Gopher" {
		t.Fatalf("内容が一致しません。実際: %+v", book)
	}
}

func TestCreateBookPersists(t *testing.T) {
	db := newTestDB(t)

	created, err := createBook(db, "テスト駆動開発", "Kent")
	if err != nil {
		t.Fatalf("createBook がエラーを返しました: %v", err)
	}

	var found Book
	if err := db.First(&found, created.ID).Error; err != nil {
		t.Fatalf("作成したレコードをDBから取得できません: %v", err)
	}
	if found.Title != "テスト駆動開発" || found.Author != "Kent" {
		t.Fatalf("保存内容が一致しません。実際: %+v", found)
	}
}

func TestListBooksEmpty(t *testing.T) {
	db := newTestDB(t)

	books, err := listBooks(db)
	if err != nil {
		t.Fatalf("listBooks がエラーを返しました: %v", err)
	}
	if len(books) != 0 {
		t.Fatalf("空のテーブルからは0件が返るはずですが、%d 件でした", len(books))
	}
}

func TestListBooksReturnsAll(t *testing.T) {
	db := newTestDB(t)

	titles := []string{"第1巻", "第2巻", "第3巻"}
	for _, title := range titles {
		if _, err := createBook(db, title, "著者"); err != nil {
			t.Fatalf("createBook がエラーを返しました: %v", err)
		}
	}

	books, err := listBooks(db)
	if err != nil {
		t.Fatalf("listBooks がエラーを返しました: %v", err)
	}
	if len(books) != len(titles) {
		t.Fatalf("件数が一致しません。期待値: %d, 実際: %d", len(titles), len(books))
	}

	seen := map[string]bool{}
	for _, b := range books {
		seen[b.Title] = true
	}
	for _, title := range titles {
		if !seen[title] {
			t.Fatalf("%q が結果に含まれていません。実際: %+v", title, books)
		}
	}
}
