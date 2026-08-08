# Lesson 1: HTTPサーバーを立てる

Goでは標準ライブラリの `net/http` だけで、外部パッケージなしにWebサーバーを作ることができます。

`http.HandleFunc` でパスとハンドラ関数を結びつけ、`http.ListenAndServe` でサーバーを起動します。

## 課題

`helloHandler` 関数を実装してください。この関数は、リクエストが来たときにレスポンスボディとして次の文字列をそのまま書き込む必要があります。

```
Hello, CodeForge!
```

## ヒント

`http.ResponseWriter` は `io.Writer` を満たしているので、`fmt.Fprint(w, "...")` や `w.Write([]byte("..."))` で書き込めます。
