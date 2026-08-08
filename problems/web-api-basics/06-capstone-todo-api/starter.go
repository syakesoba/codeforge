package main

import (
	"net/http"
	"sync"
)

// Todo は1件のTODOを表します。
type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

var (
	mu     sync.Mutex
	todos  = map[int]Todo{}
	nextID = 1
)

// listTodosHandler は登録されている全てのTODOをJSON配列で返します。
func listTodosHandler(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. mu.Lock() / defer mu.Unlock() で保護する
	// 2. todos の中身をスライスに詰め直す
	// 3. Content-Type を application/json に設定してJSONで返す（ステータスコードは自動で200）
}

// createTodoHandler はリクエストボディ {"text": "..."} から新しいTODOを作成します。
func createTodoHandler(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. リクエストボディを {"text": string} としてデコードする
	// 2. mu.Lock() / defer mu.Unlock() で保護しつつ、nextID を採番してtodosに登録する（Done: false）
	// 3. 登録が終わったら nextID をインクリメントする
	// 4. ステータスコード201で、作成したTodoをJSONで返す
}

// updateTodoHandler はパスパラメータ {id} のTODOを更新します。
func updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. パスパラメータ "id" を strconv.Atoi で数値に変換する
	// 2. リクエストボディ {"text": string, "done": bool} をデコードする
	// 3. mu.Lock() / defer mu.Unlock() で保護しつつ、該当するTODOが存在すれば内容を上書きする
	// 4. 存在すればステータスコード200で更新後のTodoを返す
	// 5. 存在しなければステータスコード404を返す
}

// deleteTodoHandler はパスパラメータ {id} のTODOを削除します。
func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. パスパラメータ "id" を strconv.Atoi で数値に変換する
	// 2. mu.Lock() / defer mu.Unlock() で保護しつつ、該当するTODOが存在すれば削除する
	// 3. 存在すればステータスコード204を返す（本文なし）
	// 4. 存在しなければステータスコード404を返す
}

// newMux はこのレッスンで使うルーティングを登録したServeMuxを返します。
func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos", listTodosHandler)
	mux.HandleFunc("POST /todos", createTodoHandler)
	mux.HandleFunc("PUT /todos/{id}", updateTodoHandler)
	mux.HandleFunc("DELETE /todos/{id}", deleteTodoHandler)
	return mux
}

func main() {
	http.ListenAndServe(":8080", newMux())
}
