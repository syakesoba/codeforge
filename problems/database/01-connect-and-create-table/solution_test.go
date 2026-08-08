package main

import "testing"

func TestSetupDBReturnsUsableDB(t *testing.T) {
	db, err := setupDB()
	if err != nil {
		t.Fatalf("setupDB がエラーを返しました: %v", err)
	}
	if db == nil {
		t.Fatal("setupDB が nil を返しました。sql.Open の結果を返していますか？")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("DBに接続できません: %v", err)
	}
}

func TestBooksTableExists(t *testing.T) {
	db, err := setupDB()
	if err != nil {
		t.Fatalf("setupDB がエラーを返しました: %v", err)
	}
	defer db.Close()

	// books テーブルが作られていれば、INSERT/SELECT が通る
	if _, err := db.Exec(`INSERT INTO books (title, author) VALUES (?, ?)`, "Go入門", "Gopher"); err != nil {
		t.Fatalf("books テーブルへのINSERTに失敗しました。テーブル定義を確認してください: %v", err)
	}

	var title, author string
	if err := db.QueryRow(`SELECT title, author FROM books WHERE id = 1`).Scan(&title, &author); err != nil {
		t.Fatalf("books テーブルからのSELECTに失敗しました: %v", err)
	}
	if title != "Go入門" || author != "Gopher" {
		t.Fatalf("保存された内容が一致しません。title=%q, author=%q", title, author)
	}
}

func TestBooksTableHasAutoIncrementID(t *testing.T) {
	db, err := setupDB()
	if err != nil {
		t.Fatalf("setupDB がエラーを返しました: %v", err)
	}
	defer db.Close()

	for i := 0; i < 3; i++ {
		if _, err := db.Exec(`INSERT INTO books (title, author) VALUES (?, ?)`, "t", "a"); err != nil {
			t.Fatalf("INSERTに失敗しました: %v", err)
		}
	}

	var maxID int
	if err := db.QueryRow(`SELECT MAX(id) FROM books`).Scan(&maxID); err != nil {
		t.Fatalf("SELECTに失敗しました: %v", err)
	}
	if maxID != 3 {
		t.Fatalf("idが自動採番されていないようです。3件INSERTした後の最大idは3のはずですが、実際は %d でした", maxID)
	}
}
