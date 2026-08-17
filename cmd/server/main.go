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
	"github.com/syakesoba/codeforge/internal/lint"
	"github.com/syakesoba/codeforge/internal/problems"
	"github.com/syakesoba/codeforge/internal/ratelimit"
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

type checkResponse struct {
	Diagnostics []lint.Diagnostic `json:"diagnostics"`
}

type formatRequest struct {
	Code string `json:"code"`
	// Hints は呼び出し側（/checkの結果）で既に判明している未解決の識別子名。
	// 指定するとgo buildによる再検証を省略できるため応答が速くなる。
	Hints []string `json:"hints"`
}

type formatResponse struct {
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

	// ログイン試行のブルートフォース対策。IPアドレス単位で制限する
	// （単一プロセスでの運用のためインメモリ実装。リバースプロキシ配下での
	// 運用に切り替える場合はクライアントIPの取得方法を見直すこと）。
	loginLimiter := ratelimit.New(10, time.Minute)
	signupLimiter := ratelimit.New(5, time.Minute)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/problems/{id}", getProblemHandler())
	mux.HandleFunc("GET /api/problems/{id}/answer", getAnswerHandler())
	mux.HandleFunc("POST /api/problems/{id}/submit", requireCSRF(submitHandler(runner, authSvc, st)))
	mux.HandleFunc("POST /api/problems/{id}/check", checkHandler())
	mux.HandleFunc("POST /api/problems/{id}/format", formatHandler())
	mux.HandleFunc("GET /api/problems/{id}/draft", getDraftHandler(authSvc, st))
	mux.HandleFunc("POST /api/problems/{id}/draft", requireCSRF(saveDraftHandler(authSvc, st)))

	mux.HandleFunc("POST /api/auth/signup", rateLimited(signupLimiter, requireCSRF(signUpHandler(authSvc, secureCookie))))
	mux.HandleFunc("POST /api/auth/login", rateLimited(loginLimiter, requireCSRF(logInHandler(authSvc, secureCookie))))
	mux.HandleFunc("POST /api/auth/logout", requireCSRF(logOutHandler(authSvc, secureCookie)))
	mux.HandleFunc("GET /api/me", meHandler(authSvc))
	mux.HandleFunc("GET /api/progress", progressHandler(authSvc, st))

	frontendOrigin := os.Getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:3000"
	}

	addr := ":8080"
	log.Printf("listening on %s (allowing CORS from %s, db=%s)", addr, frontendOrigin, dbPath)
	handler := withCORS(frontendOrigin, withCSRFCookie(secureCookie, mux))
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func withCORS(origin string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, "+csrfHeaderName)
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

// checkTimeout はエディタからの入力に対するチェックなので、採点(submit)より
// 大幅に短く設定する。go build はコンパイルのみで通常は1秒未満で終わる。
const checkTimeout = 10 * time.Second

// readTargetIfWritesTest は、ユーザーがテストコードを書く形式のレッスン
// （WritesTest）でのみ target.go（正しい実装）を読み込んで返す。
// 通常のレッスンでは nil を返す（lint.Check/FixImportsはnilなら通常の
// 実装コード用の挙動になる）。
func readTargetIfWritesTest(problem problems.Problem) ([]byte, error) {
	if !problem.WritesTest {
		return nil, nil
	}
	return problem.ReadTarget(problemsBaseDir)
}

func checkHandler() http.HandlerFunc {
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

		goMod, err := problem.ReadGoMod(problemsBaseDir)
		if err != nil {
			log.Printf("failed to read go.mod: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		goSum, err := problem.ReadGoSum(problemsBaseDir)
		if err != nil && !os.IsNotExist(err) {
			log.Printf("failed to read go.sum: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		target, err := readTargetIfWritesTest(problem)
		if err != nil {
			log.Printf("failed to read target.go: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), checkTimeout)
		defer cancel()

		diagnostics, err := lint.Check(ctx, goMod, goSum, req.Code, target)
		if err != nil {
			log.Printf("check error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(checkResponse{Diagnostics: diagnostics}); err != nil {
			log.Printf("failed to encode response: %v", err)
		}
	}
}

func formatHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		problem, ok := problems.Get(id)
		if !ok {
			http.Error(w, "problem not found", http.StatusNotFound)
			return
		}

		var req formatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		goMod, err := problem.ReadGoMod(problemsBaseDir)
		if err != nil {
			log.Printf("failed to read go.mod: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		goSum, err := problem.ReadGoSum(problemsBaseDir)
		if err != nil && !os.IsNotExist(err) {
			log.Printf("failed to read go.sum: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		target, err := readTargetIfWritesTest(problem)
		if err != nil {
			log.Printf("failed to read target.go: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), checkTimeout)
		defer cancel()

		formatted, err := lint.FixImports(ctx, goMod, goSum, req.Code, req.Hints, target)
		if err != nil {
			// 構文エラーがあるとgoimportsは整形できない。ユーザーへの
			// エラーメッセージとしてそのまま返す。
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(formatResponse{Code: formatted}); err != nil {
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
