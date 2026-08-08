package main

import (
	"errors"
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

func TestGormCRUD(t *testing.T) {
	db := newTestDB(t)

	// 1. 作成
	created, err := createBook(db, "Go入門", "Gopher")
	if err != nil {
		t.Fatalf("createBook がエラーを返しました: %v", err)
	}
	if created == nil || created.ID == 0 {
		t.Fatalf("createBook: IDが採番された *Book を返してください。実際: %+v", created)
	}

	// 2. 取得
	got, err := getBook(db, created.ID)
	if err != nil {
		t.Fatalf("getBook がエラーを返しました: %v", err)
	}
	if got == nil || got.Title != "Go入門" || got.Author != "Gopher" {
		t.Fatalf("getBook: 内容が一致しません。実際: %+v", got)
	}

	// 3. 更新
	updated, err := updateBook(db, created.ID, "Go実践", "Gopher2")
	if err != nil {
		t.Fatalf("updateBook がエラーを返しました: %v", err)
	}
	if updated == nil || updated.Title != "Go実践" || updated.Author != "Gopher2" {
		t.Fatalf("updateBook: 更新内容が反映されていません。実際: %+v", updated)
	}

	// 4. 更新がDBに反映されていること
	reloaded, err := getBook(db, created.ID)
	if err != nil {
		t.Fatalf("getBook がエラーを返しました: %v", err)
	}
	if reloaded.Title != "Go実践" {
		t.Fatalf("更新がDBに保存されていません。実際: %+v。db.Save(&book) を呼んでいますか？", reloaded)
	}

	// 5. 削除
	if err := deleteBook(db, created.ID); err != nil {
		t.Fatalf("deleteBook がエラーを返しました: %v", err)
	}

	// 6. 削除後は取得できないこと
	if _, err := getBook(db, created.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("削除後の取得は gorm.ErrRecordNotFound になるべきですが、実際: %v", err)
	}
}

func TestGetBookNotFound(t *testing.T) {
	db := newTestDB(t)

	_, err := getBook(db, 999999)

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("存在しないIDの取得は gorm.ErrRecordNotFound になるべきですが、実際: %v", err)
	}
}

func TestUpdateBookNotFound(t *testing.T) {
	db := newTestDB(t)

	_, err := updateBook(db, 999999, "x", "y")

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("存在しないIDの更新は gorm.ErrRecordNotFound になるべきですが、実際: %v", err)
	}
}

func TestDeleteBookNotFound(t *testing.T) {
	db := newTestDB(t)

	err := deleteBook(db, 999999)

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("存在しないIDの削除は gorm.ErrRecordNotFound になるべきですが、実際: %v。result.RowsAffected を確認していますか？", err)
	}
}

func TestDeleteBookDoesNotAffectOthers(t *testing.T) {
	db := newTestDB(t)

	keep, _ := createBook(db, "残る本", "著者A")
	remove, _ := createBook(db, "消える本", "著者B")

	if err := deleteBook(db, remove.ID); err != nil {
		t.Fatalf("deleteBook がエラーを返しました: %v", err)
	}

	if _, err := getBook(db, keep.ID); err != nil {
		t.Fatalf("他のレコードまで削除されています: %v", err)
	}
}
