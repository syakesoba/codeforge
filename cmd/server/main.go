// Command server は CodeForge の採点APIを提供するHTTPサーバーです。
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/syakesoba/codeforge/internal/auth"
	"github.com/syakesoba/codeforge/internal/judge"
	"github.com/syakesoba/codeforge/internal/problems"
	"github.com/syakesoba/codeforge/internal/store"
)

const problemsBaseDir = "problems"

// dbPath はユーザー・セッション・学習進捗を保存するSQLiteファイルのパスです。
const dbPath = "data/codeforge.db"

type submitRequest struct {
	Code string `json:"code"`
}

type problemResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Markdown    string `json:"markdown"`
	StarterCode string `json:"starterCode"`
}

type answerResponse struct {
	Code string `json:"code"`
}

func main() {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		log.Fatalf("failed to create data directory: %v", err)
	}

	st, err := store.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open store: %v", err)
	}
	defer st.Close()

	authSvc := auth.New(st)
	runner := judge.NewRunner()

	// SameSite=None は Secure 必須で、ローカル開発(http)ではCookieが拒否される。
	// 本番(HTTPS)では SECURE_COOKIE=true を設定すること。
	secureCookie := os.Getenv("SECURE_COOKIE") == "true"

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/problems/{id}", getProblemHandler())
	mux.HandleFunc("GET /api/problems/{id}/answer", getAnswerHandler())
	mux.HandleFunc("POST /api/problems/{id}/submit", submitHandler(runner, authSvc, st))

	mux.HandleFunc("POST /api/auth/signup", signUpHandler(authSvc, secureCookie))
	mux.HandleFunc("POST /api/auth/login", logInHandler(authSvc, secureCookie))
	mux.HandleFunc("POST /api/auth/logout", logOutHandler(authSvc, secureCookie))
	mux.HandleFunc("GET /api/me", meHandler(authSvc))
	mux.HandleFunc("GET /api/progress", progressHandler(authSvc, st))

	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:3000"
	}

	addr := ":8080"
	log.Printf("listening on %s (allowing CORS from %s, db=%s)", addr, frontendOrigin, dbPath)
	if err := http.ListenAndServe(addr, withCORS(frontendOrigin, mux)); err != nil {
		log.Fatal(err)
	}
}

func withCORS(origin string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// Cookieによるセッションをクロスオリジンで送受信するために必要。
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func getProblemHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		problem, ok := problems.Get(id)
		if !ok {
			http.Error(w, "problem not found", http.StatusNotFound)
			return
		}

		markdown, err := problem.ReadMarkdown(problemsBaseDir)
		if err != nil {
			log.Printf("failed to read markdown: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		starter, err := problem.ReadStarter(problemsBaseDir)
		if err != nil {
			log.Printf("failed to read starter code: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		resp := problemResponse{
			ID:          problem.ID,
			Title:       problem.Title,
			Markdown:    string(markdown),
			StarterCode: string(starter),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("failed to encode response: %v", err)
		}
	}
}

// stripBuildIgnoreTag は answer.go 先頭の `//go:build ignore` 行
// （ビルド対象から除外するための実装上のタグ）を、表示用に取り除く。
func stripBuildIgnoreTag(code string) string {
	const tag = "//go:build ignore\n\n"
	return strings.TrimPrefix(code, tag)
}

func getAnswerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		problem, ok := problems.Get(id)
		if !ok {
			http.Error(w, "problem not found", http.StatusNotFound)
			return
		}

		answer, err := problem.ReadAnswer(problemsBaseDir)
		if err != nil {
			log.Printf("failed to read answer code: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(answerResponse{Code: stripBuildIgnoreTag(string(answer))}); err != nil {
			log.Printf("failed to encode response: %v", err)
		}
	}
}

func submitHandler(runner judge.Runner, authSvc *auth.Service, st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		problem, ok := problems.Get(id)
		if !ok {
			http.Error(w, "problem not found", http.StatusNotFound)
			return
		}

		var req submitRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// judge側のタイムアウト（30秒）より長くしておき、
		// サンドボックス側のタイムアウト処理が先に働くようにする。
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		result, err := runner.Run(ctx, problem, req.Code)
		if err != nil {
			log.Printf("judge run error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// ログイン中なら合格を記録する。未ログインでも採点自体は利用できる
		// （その場合、進捗はフロントエンドのlocalStorageにのみ保存される）。
		if result.Passed {
			if user, err := authSvc.UserByToken(sessionToken(r)); err == nil {
				if err := st.MarkProgress(user.ID, problem.ID); err != nil {
					log.Printf("failed to mark progress: %v", err)
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			log.Printf("failed to encode response: %v", err)
		}
	}
}
