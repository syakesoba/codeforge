//go:build ignore

package main

import (
	"encoding/json"
	"net/http"
	"strconv"
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
	mu.Lock()
	defer mu.Unlock()

	result := make([]Todo, 0, len(todos))
	for _, t := range todos {
		result = append(result, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

type createTodoRequest struct {
	Text string `json:"text"`
}

// createTodoHandler はリクエストボディ {"text": "..."} から新しいTODOを作成します。
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

type updateTodoRequest struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

// updateTodoHandler はパスパラメータ {id} のTODOを更新します。
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

// deleteTodoHandler はパスパラメータ {id} のTODOを削除します。
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
