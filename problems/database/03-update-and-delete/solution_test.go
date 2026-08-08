package main

import (
	"database/sql"
	"errors"
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

func TestUpdateBook(t *testing.T) {
	db := newTestDB(t)

	id, err := insertBook(db, "旧タイトル", "旧著者")
	if err != nil {
		t.Fatalf("insertBook がエラーを返しました: %v", err)
	}

	if err := updateBook(db, id, "新タイトル", "新著者"); err != nil {
		t.Fatalf("updateBook がエラーを返しました: %v", err)
	}

	var title, author string
	if err := db.QueryRow(`SELECT title, author FROM books WHERE id = ?`, id).Scan(&title, &author); err != nil {
		t.Fatalf("SELECTに失敗しました: %v", err)
	}
	if title != "新タイトル" || author != "新著者" {
		t.Fatalf("更新内容が反映されていません。title=%q, author=%q", title, author)
	}
}

func TestUpdateBookNotFound(t *testing.T) {
	db := newTestDB(t)

	err := updateBook(db, 999999, "x", "y")

	if err == nil {
		t.Fatal("存在しないIDの更新はエラーになるべきですが、nilが返りました。RowsAffected() を確認していますか？")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("ErrNotFound が返るべきですが、実際は: %v", err)
	}
}

func TestDeleteBook(t *testing.T) {
	db := newTestDB(t)

	id, err := insertBook(db, "消される本", "著者")
	if err != nil {
		t.Fatalf("insertBook がエラーを返しました: %v", err)
	}

	if err := deleteBook(db, id); err != nil {
		t.Fatalf("deleteBook がエラーを返しました: %v", err)
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM books WHERE id = ?`, id).Scan(&count); err != nil {
		t.Fatalf("SELECTに失敗しました: %v", err)
	}
	if count != 0 {
		t.Fatalf("削除されていません。該当する行が %d 件残っています", count)
	}
}

func TestDeleteBookNotFound(t *testing.T) {
	db := newTestDB(t)

	err := deleteBook(db, 999999)

	if err == nil {
		t.Fatal("存在しないIDの削除はエラーになるべきですが、nilが返りました。RowsAffected() を確認していますか？")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("ErrNotFound が返るべきですが、実際は: %v", err)
	}
}

func TestDeleteBookDoesNotAffectOthers(t *testing.T) {
	db := newTestDB(t)

	id1, _ := insertBook(db, "残る本", "著者A")
	id2, _ := insertBook(db, "消える本", "著者B")

	if err := deleteBook(db, id2); err != nil {
		t.Fatalf("deleteBook がエラーを返しました: %v", err)
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM books WHERE id = ?`, id1).Scan(&count); err != nil {
		t.Fatalf("SELECTに失敗しました: %v", err)
	}
	if count != 1 {
		t.Fatal("他の行まで削除されています。WHERE句は付いていますか？")
	}
}
