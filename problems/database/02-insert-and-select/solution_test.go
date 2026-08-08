package main

import (
	"database/sql"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := setupDB()
	if err != nil {
		t.Fatalf("setupDB がエラーを返しました: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestInsertBookReturnsID(t *testing.T) {
	db := newTestDB(t)

	id1, err := insertBook(db, "Go入門", "Gopher")
	if err != nil {
		t.Fatalf("insertBook がエラーを返しました: %v", err)
	}
	if id1 != 1 {
		t.Fatalf("1件目のIDは1のはずですが、実際は %d でした。LastInsertId() を返していますか？", id1)
	}

	id2, err := insertBook(db, "Go実践", "Gopher2")
	if err != nil {
		t.Fatalf("insertBook がエラーを返しました: %v", err)
	}
	if id2 != 2 {
		t.Fatalf("2件目のIDは2のはずですが、実際は %d でした", id2)
	}
}

func TestInsertBookStoresValues(t *testing.T) {
	db := newTestDB(t)

	if _, err := insertBook(db, "テスト駆動開発", "Kent"); err != nil {
		t.Fatalf("insertBook がエラーを返しました: %v", err)
	}

	var title, author string
	if err := db.QueryRow(`SELECT title, author FROM books WHERE id = 1`).Scan(&title, &author); err != nil {
		t.Fatalf("SELECTに失敗しました: %v", err)
	}
	if title != "テスト駆動開発" || author != "Kent" {
		t.Fatalf("保存された内容が一致しません。title=%q, author=%q", title, author)
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

func TestListBooksReturnsAllInIDOrder(t *testing.T) {
	db := newTestDB(t)

	want := []Book{
		{ID: 1, Title: "第1巻", Author: "著者A"},
		{ID: 2, Title: "第2巻", Author: "著者B"},
		{ID: 3, Title: "第3巻", Author: "著者C"},
	}
	for _, b := range want {
		if _, err := insertBook(db, b.Title, b.Author); err != nil {
			t.Fatalf("insertBook がエラーを返しました: %v", err)
		}
	}

	got, err := listBooks(db)
	if err != nil {
		t.Fatalf("listBooks がエラーを返しました: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("件数が一致しません。期待値: %d, 実際: %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%d番目の要素が一致しません。期待値: %+v, 実際: %+v（id順に並んでいますか？）", i, want[i], got[i])
		}
	}
}
