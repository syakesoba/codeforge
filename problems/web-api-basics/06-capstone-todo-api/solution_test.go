package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTodoAPI(t *testing.T) {
	mux := newMux()

	// 1. 初期状態は空であること
	{
		req := httptest.NewRequest(http.MethodGet, "/todos", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /todos: 期待したステータスコードは %d でしたが、実際は %d でした", http.StatusOK, rec.Code)
		}
		var got []Todo
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("GET /todos: レスポンスのJSONデコードに失敗しました: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("GET /todos: 初期状態は0件であるべきですが %d 件でした", len(got))
		}
	}

	// 2. 作成できること
	var created Todo
	{
		req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(`{"text":"牛乳を買う"}`))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("POST /todos: 期待したステータスコードは %d でしたが、実際は %d でした", http.StatusCreated, rec.Code)
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
			t.Fatalf("POST /todos: レスポンスのJSONデコードに失敗しました: %v", err)
		}
		if created.Text != "牛乳を買う" {
			t.Fatalf("POST /todos: textが一致しません。期待値: %q, 実際: %q", "牛乳を買う", created.Text)
		}
		if created.Done {
			t.Fatalf("POST /todos: 作成直後のdoneはfalseであるべきです")
		}
		if created.ID == 0 {
			t.Fatalf("POST /todos: idが採番されていません")
		}
	}

	// 3. 一覧に反映されていること
	{
		req := httptest.NewRequest(http.MethodGet, "/todos", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		var got []Todo
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("GET /todos: レスポンスのJSONデコードに失敗しました: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("GET /todos: 作成後は1件であるべきですが %d 件でした", len(got))
		}
	}

	// 4. 更新できること
	{
		path := fmt.Sprintf("/todos/%d", created.ID)
		req := httptest.NewRequest(http.MethodPut, path, strings.NewReader(`{"text":"牛乳を買う","done":true}`))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("PUT %s: 期待したステータスコードは %d でしたが、実際は %d でした", path, http.StatusOK, rec.Code)
		}
		var updated Todo
		if err := json.Unmarshal(rec.Body.Bytes(), &updated); err != nil {
			t.Fatalf("PUT %s: レスポンスのJSONデコードに失敗しました: %v", path, err)
		}
		if !updated.Done {
			t.Fatalf("PUT %s: doneがtrueに更新されていません", path)
		}
	}

	// 5. 存在しないIDの更新は404であること
	{
		req := httptest.NewRequest(http.MethodPut, "/todos/999999", strings.NewReader(`{"text":"x","done":false}`))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("PUT /todos/999999: 期待したステータスコードは %d でしたが、実際は %d でした", http.StatusNotFound, rec.Code)
		}
	}

	// 6. 削除できること
	{
		path := fmt.Sprintf("/todos/%d", created.ID)
		req := httptest.NewRequest(http.MethodDelete, path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Fatalf("DELETE %s: 期待したステータスコードは %d でしたが、実際は %d でした", path, http.StatusNoContent, rec.Code)
		}
	}

	// 7. 削除後は一覧が空であること
	{
		req := httptest.NewRequest(http.MethodGet, "/todos", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		var got []Todo
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("GET /todos: レスポンスのJSONデコードに失敗しました: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("GET /todos: 削除後は0件であるべきですが %d 件でした", len(got))
		}
	}

	// 8. 存在しないIDの削除は404であること
	{
		req := httptest.NewRequest(http.MethodDelete, "/todos/999999", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("DELETE /todos/999999: 期待したステータスコードは %d でしたが、実際は %d でした", http.StatusNotFound, rec.Code)
		}
	}
}
