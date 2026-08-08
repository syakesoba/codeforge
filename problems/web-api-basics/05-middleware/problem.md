# Lesson 5: 共通処理をまとめる（ミドルウェア）

## 今回学ぶこと

- 複数のハンドラに共通する処理を1箇所にまとめる「ミドルウェア」というパターン
- `http.Handler` を受け取って `http.Handler` を返す関数の書き方

## 解説

「アクセスログを出力する」「レスポンスに共通ヘッダーを付ける」「認証をチェックする」など、多くのエンドポイントに共通する処理は、ハンドラ1つ1つに書くのではなく「ミドルウェア」としてまとめるのがGoでの定石です。

ミドルウェアの正体は、**「`http.Handler` を受け取って、別の `http.Handler` を返す関数」**です。少し抽象的なので、実際のコードで見てみましょう。次の例は、リクエストが来るたびにログを出力するミドルウェアです。

```go
func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r) // 元のハンドラを呼び出す
	})
}
```

ポイントは3つです。

1. 引数 `next` は「本来呼ばれるはずだったハンドラ」
2. 戻り値は `http.HandlerFunc(func(w, r) {...})` という形の**新しいハンドラ**
3. 新しいハンドラの中で、共通処理（ここではログ出力）をしたあとに `next.ServeHTTP(w, r)` を呼んで元の処理につなげる

使うときは、既存のハンドラ（や `*http.ServeMux`）を関数でくるむだけです。

```go
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", helloHandler)

	http.ListenAndServe(":8080", withLogging(mux))
}
```

こうすることで、`helloHandler` 自体には手を加えずに、すべてのリクエストにログ出力を追加できました。同じ形で「共通ヘッダーを付ける」「認証をチェックする」といった処理も作れます。

## 演習

`withPoweredByHeader` 関数を実装してください。この関数はミドルウェアで、次の処理を行う必要があります。

1. レスポンスヘッダー `X-Powered-By` に `CodeForge` を設定する
2. その後、`next.ServeHTTP(w, r)` を呼び出して元のハンドラの処理につなげる

つまり、`withPoweredByHeader` でラップされたどんなハンドラも、レスポンスに `X-Powered-By: CodeForge` ヘッダーが付くようになります。

## ヒント

- 型は `func withPoweredByHeader(next http.Handler) http.Handler` とすでに定義されています
- ヘッダーの設定は本文を書き込む**前**に行う必要があります（Lesson 4を参照）
