# Lesson 4: レスポンスを返す

## 今回学ぶこと

- JSON形式でレスポンスを返す方法（`json.NewEncoder(w).Encode(...)`）
- レスポンスヘッダー・ステータスコードの設定方法
- ヘッダー・ステータスコード・ボディを書く順番の重要な注意点

## 解説

### JSONを返す

これまでは `fmt.Fprint` で文字列を返してきましたが、実際のAPIではJSON形式でレスポンスを返すことが一般的です。`encoding/json` パッケージの `json.NewEncoder(w).Encode(値)` を使うと、Goの構造体をJSONに変換してそのまま書き込めます。

```go
type messageResponse struct {
	Message string `json:"message"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messageResponse{Message: "pong"})
}
```

`w.Header().Set("Content-Type", "application/json")` を書いておくことで、クライアント（ブラウザやAPIクライアント）に「このレスポンスはJSONですよ」と伝えられます。省略しても動作はしますが、実務では必ず設定する習慣をつけましょう。

### ステータスコードを明示的に返す

何も指定しない場合、Goは自動的にステータスコード `200 OK` を返します。作成系のAPI（何かを新規作成するAPI）では `201 Created` を返すのが一般的です。明示的にステータスコードを指定するには `w.WriteHeader(コード)` を使います。

```go
func createPingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(messageResponse{Message: "created"})
}
```

> **重要な注意点: 呼び出す順番**
>
> `w.Header().Set(...)` → `w.WriteHeader(...)` → 本文の書き込み（`json.NewEncoder(w).Encode(...)` や `w.Write(...)`）の順番を必ず守ってください。
>
> `http.ResponseWriter` は、最初に何か書き込まれた時点（`Write` が最初に呼ばれた瞬間）で「ステータスコード200・それまでに設定されたヘッダー」が確定してしまいます。つまり、本文を書き込んだ**あとに** `w.WriteHeader` や `w.Header().Set` を呼んでも、レスポンスにはもう反映されません。「ヘッダー設定 → ステータスコード確定 → 本文」の順番、と覚えておきましょう。

## 演習

`createGreetingHandler` 関数を実装してください。この関数は次の3つを行う必要があります。

1. レスポンスヘッダー `Content-Type` に `application/json` を設定する
2. ステータスコードとして `201 Created` を返す
3. レスポンスボディとして次のJSONを返す

```json
{"message": "Hello, CodeForge!"}
```

レスポンス用のstruct `greetingResponse` はすでに用意されています。

## ヒント

- 順番は `w.Header().Set(...)` → `w.WriteHeader(http.StatusCreated)` → `json.NewEncoder(w).Encode(...)` です
