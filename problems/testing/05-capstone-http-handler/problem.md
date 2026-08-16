# Lesson 5（道場）: HTTPハンドラーをテストする

## 今回学ぶこと

- `net/http/httptest` を使い、実際にサーバーを起動せずHTTPハンドラーをテストする方法
- `httptest.NewRequest` / `httptest.NewRecorder` の使い方
- これまで学んだテーブル駆動テスト・サブテスト・エラーケースの総仕上げ

### httptest で「サーバーを起動せずに」テストする

「Web APIの基礎」コースで作ったハンドラーは `http.HandlerFunc` 型の関数でした。これは「`http.ResponseWriter` と `*http.Request` を受け取る普通の関数」なので、実際にサーバーを起動してHTTPリクエストを送らなくても、**直接関数として呼び出す**ことでテストできます。

```go
func greetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello!")
}

func TestGreetHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/greet", nil)
	rec := httptest.NewRecorder()

	greetHandler(rec, req) // サーバー無しで直接呼び出す

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "Hello!" {
		t.Errorf("body = %q, want %q", rec.Body.String(), "Hello!")
	}
}
```

- `httptest.NewRequest(method, target, body)`: テスト用の `*http.Request` を作る。`target` にはパスやクエリ文字列（`"/greet?name=Alice"`）を書ける
- `httptest.NewRecorder()`: `http.ResponseWriter` を実装した「記録係」を作る。ハンドラーが書き込んだステータスコードや本文を、後から `rec.Code` / `rec.Body.String()` で確認できる

### これまで学んだことの総仕上げ

このレッスンは「道場」＝総仕上げ問題です。Lesson 1〜4で学んだ以下の要素を、すべて組み合わせて使います。

- テーブル駆動テスト（Lesson 1）: 複数のケースを1つのループで検証する
- サブテスト（Lesson 2）: `t.Run` でケースごとに名前を付けて実行する
- 「正常系」と「異常系（エラー）」を1つのテーブルで扱う考え方（Lesson 4）: 今回は `wantStatus` でステータスコードによって検証内容を分岐させます

## 演習

以下の `greetHandler` は、クエリパラメータ `name` を読み取って挨拶文を返すハンドラーで、**すでに実装済み**です（変更不要です）。`name` が空の場合は `400 Bad Request` を返します。

```go
func greetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Hello, %s!", name)
}
```

`runGreetHandlerCases` を実装してください。`cases` の各要素について、`httptest.NewRequest` と `httptest.NewRecorder` でリクエスト/レスポンスを用意し、`greetHandler` を直接呼び出して、ステータスコード（`rec.Code`）が `wantStatus` と一致するか確認してください。ステータスコードが `http.StatusOK`（200）のケースについては、本文（`rec.Body.String()`）が `wantBody` と一致するかも確認してください。

## ヒント

- リクエストの作り方: `httptest.NewRequest(http.MethodGet, "/greet"+c.query, nil)`（`c.query` は `"?name=Alice"` のような文字列）
- レコーダーの作り方: `httptest.NewRecorder()`
- ハンドラーの呼び出し方: `greetHandler(rec, req)`（他の関数と同じように直接呼べます）
- ステータスコードの比較には `rec.Code`、本文の比較には `rec.Body.String()` を使います
