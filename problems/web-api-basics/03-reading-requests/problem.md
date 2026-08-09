# Lesson 3: リクエストを読み取る

## 今回学ぶこと

- クエリパラメータの読み取り方（`r.URL.Query().Get(...)`）
- リクエストボディ（JSON）のデコード方法（`json.NewDecoder(r.Body).Decode(...)`）

### クエリパラメータを読み取る

`/search?q=golang` のようなURLの `?` 以降の部分を「クエリパラメータ」と呼びます。Goでは `r.URL.Query()` で全パラメータを取得し、`.Get("キー名")` で値を取り出せます。

```go
func searchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	fmt.Fprintf(w, "You searched: %s", q)
}
```

`/search?q=golang` にアクセスすると `You searched: golang` が返ります。存在しないキーを指定した場合は空文字列 `""` が返ります。

### JSONのリクエストボディを読み取る

APIでは、クライアントがJSON形式のデータをリクエストボディに入れて送ってくることがよくあります。Goでは対応するstruct（構造体）を用意し、`encoding/json` パッケージの `json.NewDecoder(r.Body).Decode(&変数)` でデコードします。

```go
type greetRequest struct {
	Name string `json:"name"`
}

func greetJSONHandler(w http.ResponseWriter, r *http.Request) {
	var req greetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Hello, %s!", req.Name)
}
```

structのフィールドに書いてある `` `json:"name"` `` の部分は「構造体タグ」と呼ばれ、JSONのキー `"name"` とGoのフィールド `Name` を対応づけています。例えば `{"name": "Gopher"}` というリクエストボディが送られてくると、`req.Name` に `"Gopher"` が入ります。

デコードに失敗した場合（JSONの形式が壊れているなど）は `err` が `nil` ではなくなるので、`http.Error` でエラーレスポンスを返すのが定石です。

## 演習

`addHandler` 関数を実装してください。この関数は、リクエストボディとして次のようなJSONを受け取ります。

```json
{"a": 3, "b": 4}
```

`a` と `b` を足し算し、その結果を**文字列として**レスポンスボディに書き込んでください（上の例なら `7` という文字列）。

デコード用のstruct `addRequest` はすでに用意されています。JSONのキー `"a"` は `A` フィールドに、`"b"` は `B` フィールドに対応しています。

## ヒント

- `json.NewDecoder(r.Body).Decode(&req)` でリクエストボディをデコードできます
- 合計値は `int` なので、文字列に変換するには `fmt.Fprintf(w, "%d", sum)` が使えます
