# Lesson 6（道場）: TODO管理APIを作ろう

## 今回学ぶこと

このレッスンは「道場」、つまりCourse「Web APIの基礎」のまとめ問題です。Lesson 1〜5で学んだ内容をすべて組み合わせて、TODO（やることリスト）を管理するREST APIを作ります。

- Lesson 1: ハンドラの基本
- Lesson 2: ルーティングとパスパラメータ
- Lesson 3: JSONリクエストボディの読み取り
- Lesson 4: JSONレスポンス・ステータスコード
- Lesson 5（今回は直接使いませんが）: 共通処理のまとめ方

## 解説: 作るものの仕様

メモリ上（変数）でTODOを管理する、次の4つのエンドポイントを持つAPIを作ります。

| メソッド | パス | 役割 |
|---|---|---|
| `GET` | `/todos` | 登録されている全TODOを一覧で返す |
| `POST` | `/todos` | 新しいTODOを作成する |
| `PUT` | `/todos/{id}` | 指定したIDのTODOを更新する |
| `DELETE` | `/todos/{id}` | 指定したIDのTODOを削除する |

TODO 1件のデータ構造（`Todo` 型）はすでに定義されています。

```go
type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}
```

保存先として、パッケージレベルの `map[int]Todo` と、次のIDを管理する `nextID` もすでに用意されています（`sync.Mutex` で保護しているので、そのまま使ってください）。

```go
var (
	mu     sync.Mutex
	todos  = map[int]Todo{}
	nextID = 1
)
```

複数のリクエストが同時に来ても壊れないように、`todos` や `nextID` を読み書きする前後は必ず `mu.Lock()` / `defer mu.Unlock()` で保護してください（Lesson 1〜5では扱いませんでしたが、共有の変数を複数のハンドラから触るときの定石です）。

### listTodosHandler（一覧取得）

`todos` に入っている全TODOをJSONの配列として返します。ステータスコードは `200`（何も指定しなければ自動的に200になります）。

```go
func listTodosHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	result := make([]Todo, 0, len(todos))
	for _, t := range todos {
		result = append(result, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
```

### createTodoHandler（作成）

リクエストボディ `{"text": "牛乳を買う"}` を受け取り、`nextID` を割り当てて `todos` に保存します。`Done` は `false` で初期化します。保存が終わったら `nextID` をインクリメントするのを忘れないでください。レスポンスはステータスコード `201`（`http.StatusCreated`）で、作成したTodoをJSONで返します。

```go
type createTodoRequest struct {
	Text string `json:"text"`
}

func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	mu.Lock()
	todo := Todo{ID: nextID, Text: req.Text, Done: false}
	todos[nextID] = todo
	nextID++
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}
```

`createTodoRequest` はこのハンドラの中だけで使う型なので、`starter.go` には定義されていません。自分で追加してください（Lesson 3の `addRequest` と同じ考え方です）。

### パスパラメータを数値に変換する

`updateTodoHandler` と `deleteTodoHandler` はどちらもパスパラメータ `{id}` を使います。Lesson 2で学んだ `r.PathValue("id")` は文字列を返すので、数値として使うには変換が必要です。`strconv.Atoi` を使います。

```go
id, err := strconv.Atoi(r.PathValue("id"))
if err != nil {
	http.Error(w, "invalid id", http.StatusBadRequest)
	return
}
```

### updateTodoHandler（更新）

パスパラメータ `{id}` を数値に変換し、該当するTODOが存在すれば、リクエストボディ `{"text": "...", "done": true}` の内容で上書きします。存在すればステータスコード `200` で更新後のTodoを返し、存在しなければステータスコード `404`（`http.StatusNotFound`）を返してください（`http.Error` を使えば本文もエラーメッセージも自動で設定されます）。

```go
type updateTodoRequest struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req updateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	todo, ok := todos[id]
	if !ok {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	todo.Text = req.Text
	todo.Done = req.Done
	todos[id] = todo

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}
```

### deleteTodoHandler（削除）

パスパラメータ `{id}` を数値に変換し、該当するTODOが存在すれば `delete(todos, id)` で削除してステータスコード `204`（`http.StatusNoContent`、本文なし）を返します。存在しなければ `404` を返してください。

```go
func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, ok := todos[id]; !ok {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}
	delete(todos, id)

	w.WriteHeader(http.StatusNoContent)
}
```

204（No Content）は本文を返さないステータスコードなので、`json.NewEncoder(w).Encode(...)` は呼ばない点に注意してください。

### ルーティングの登録

`newMux()` の中で、4つのエンドポイントはすでに登録済みです。あなたが実装するのは4つのハンドラ関数の中身だけです。

```go
func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", listTodosHandler)
	mux.HandleFunc("POST /todos", createTodoHandler)
	mux.HandleFunc("PUT /todos/{id}", updateTodoHandler)
	mux.HandleFunc("DELETE /todos/{id}", deleteTodoHandler)
	return mux
}
```

## 演習

上記の仕様どおりに、4つのハンドラ関数（`listTodosHandler`, `createTodoHandler`, `updateTodoHandler`, `deleteTodoHandler`）を実装してください。

採点では、作成 → 一覧確認 → 更新 → 存在しないIDへの更新（404確認）→ 削除 → 一覧確認 → 存在しないIDへの削除（404確認）という一連の流れを自動でテストします。

## ヒント

- レスポンスを返す前に必ず `mu.Lock()` / `defer mu.Unlock()` する
- ステータスコードとヘッダーを設定する順番は Lesson 4 のとおり「ヘッダー → ステータスコード → 本文」
- 存在しないIDの場合は `http.Error(w, "todo not found", http.StatusNotFound)` のように書けます
