package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func doRequest(t *testing.T, r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestBooksAPI(t *testing.T) {
	r := newRouter()

	// 1. 初期状態は空であること
	{
		rec := doRequest(t, r, http.MethodGet, "/api/books", "")
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /api/books: 期待したステータスコードは %d でしたが、実際は %d でした", http.StatusOK, rec.Code)
		}
		var got []Book
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("GET /api/books: JSONデコードに失敗しました: %v (body=%q)", err, rec.Body.String())
		}
		if len(got) != 0 {
			t.Fatalf("GET /api/books: 初期状態は0件であるべきですが %d 件でした", len(got))
		}
	}

	// 2. 作成できること
	var created Book
	{
		rec := doRequest(t, r, http.MethodPost, "/api/books", `{"title":"Go入門","author":"Gopher"}`)
		if rec.Code != http.StatusCreated {
			t.Fatalf("POST /api/books: 期待したステータスコードは %d でしたが、実際は %d でした (body=%q)",
				http.StatusCreated, rec.Code, rec.Body.String())
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
			t.Fatalf("POST /api/books: JSONデコードに失敗しました: %v (body=%q)", err, rec.Body.String())
		}
		if created.Title != "Go入門" || created.Author != "Gopher" {
			t.Fatalf("POST /api/books: 内容が一致しません。実際: %+v", created)
		}
		if created.ID == 0 {
			t.Fatal("POST /api/books: idが採番されていません")
		}
	}

	// 3. バリデーションエラーは400であること
	{
		rec := doRequest(t, r, http.MethodPost, "/api/books", `{"title":"タイトルだけ"}`)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("POST /api/books (authorなし): 期待したステータスコードは %d でしたが、実際は %d でした。binding:\"required\" は付いていますか？",
				http.StatusBadRequest, rec.Code)
		}
	}

	// 4. 一覧に反映されていること
	{
		rec := doRequest(t, r, http.MethodGet, "/api/books", "")
		var got []Book
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("GET /api/books: JSONデコードに失敗しました: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("GET /api/books: 作成後は1件であるべきですが %d 件でした", len(got))
		}
	}

	// 5. 更新できること
	{
		path := fmt.Sprintf("/api/books/%d", created.ID)
		rec := doRequest(t, r, http.MethodPut, path, `{"title":"Go実践","author":"Gopher2"}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("PUT %s: 期待したステータスコードは %d でしたが、実際は %d でした (body=%q)",
				path, http.StatusOK, rec.Code, rec.Body.String())
		}
		var updated Book
		if err := json.Unmarshal(rec.Body.Bytes(), &updated); err != nil {
			t.Fatalf("PUT %s: JSONデコードに失敗しました: %v", path, err)
		}
		if updated.Title != "Go実践" || updated.Author != "Gopher2" {
			t.Fatalf("PUT %s: 更新内容が反映されていません。実際: %+v", path, updated)
		}
	}

	// 6. 存在しないIDの更新は404であること
	{
		rec := doRequest(t, r, http.MethodPut, "/api/books/999999", `{"title":"x","author":"y"}`)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("PUT /api/books/999999: 期待したステータスコードは %d でしたが、実際は %d でした",
				http.StatusNotFound, rec.Code)
		}
	}

	// 7. 削除できること
	{
		path := fmt.Sprintf("/api/books/%d", created.ID)
		rec := doRequest(t, r, http.MethodDelete, path, "")
		if rec.Code != http.StatusNoContent {
			t.Fatalf("DELETE %s: 期待したステータスコードは %d でしたが、実際は %d でした (body=%q)",
				path, http.StatusNoContent, rec.Code, rec.Body.String())
		}
	}

	// 8. 削除後は一覧が空であること
	{
		rec := doRequest(t, r, http.MethodGet, "/api/books", "")
		var got []Book
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("GET /api/books: JSONデコードに失敗しました: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("GET /api/books: 削除後は0件であるべきですが %d 件でした", len(got))
		}
	}

	// 9. 存在しないIDの削除は404であること
	{
		rec := doRequest(t, r, http.MethodDelete, "/api/books/999999", "")
		if rec.Code != http.StatusNotFound {
			t.Fatalf("DELETE /api/books/999999: 期待したステータスコードは %d でしたが、実際は %d でした",
				http.StatusNotFound, rec.Code)
		}
	}
}
