# Lesson 2: ルーティングとパスパラメータ

## 今回学ぶこと

- `http.ServeMux` に複数のエンドポイントを登録する方法
- URLの一部を変数として受け取る「パスパラメータ」の使い方（`{name}` と `r.PathValue`）

## 解説

Lesson 1では `http.HandleFunc` を1つだけ使いましたが、実際のWebアプリでは複数のエンドポイント（URL）を扱います。Goでは `http.ServeMux` に複数のパターンを登録することでルーティングを組み立てます。

さらに、Go 1.22以降は `http.ServeMux` のパターンに `{名前}` という書き方でパスパラメータを埋め込めます。埋め込んだ値は、ハンドラの中で `r.PathValue("名前")` を使って取り出せます。

例えば、`/items/{id}` というパターンに対して `/items/42` というリクエストが来た場合、ハンドラの中で `r.PathValue("id")` を呼ぶと文字列 `"42"` が取得できます。

次のコードは、`/echo/{word}` というパスで受け取った単語をそのままオウム返しするサーバーの例です。

```go
package main

import (
	"fmt"
	"net/http"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	word := r.PathValue("word")
	fmt.Fprintf(w, "You said: %s", word)
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /echo/{word}", echoHandler)
	return mux
}

func main() {
	http.ListenAndServe(":8080", newMux())
}
```

ポイントは2つです。

1. `mux.HandleFunc("GET /echo/{word}", echoHandler)` のように、パターンの中に `{word}` と書くとパスパラメータになる
2. ハンドラの中で `r.PathValue("word")` を呼ぶと、リクエストされた実際の値（例: `/echo/hello` なら `"hello"`）が取れる

このレッスンからは、ルーティング設定を1箇所にまとめるために `newMux()` という関数を用意し、`main()` からも採点用のテストからもこの関数を使って `*http.ServeMux` を組み立てます。

## 演習

`greetHandler` 関数を実装してください。この関数はパスパラメータ `name` を受け取り、レスポンスボディに次の形式で文字列を書き込みます。

```
Hello, <name>!
```

例えば `/greet/Gopher` にリクエストが来た場合、レスポンスボディは `Hello, Gopher!` になります。

すでに `newMux()` の中で `GET /greet/{name}` というパターンが `greetHandler` に紐づけて登録されているので、あなたが実装するのは `greetHandler` の中身だけです。

## ヒント

- `r.PathValue("name")` で名前を取得できます
- `fmt.Fprintf(w, "Hello, %s!", name)` のように書き込めます
