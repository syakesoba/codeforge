// Command server は CodeForge の採点APIを提供するHTTPサーバーです。
package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/syakesoba/codeforge/internal/auth"
	"github.com/syakesoba/codeforge/internal/judge"
	"github.com/syakesoba/codeforge/internal/lint"
	"github.com/syakesoba/codeforge/internal/problems"
	"github.com/syakesoba/codeforge/internal/ratelimit"
	"github.com/syakesoba/codeforge/internal/store"
)

const problemsBaseDir = "problems"

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

// config はホスト環境ごとに変わる設定値です。ハードコードせず環境変数から
// 読み込むことで、ローカル・Docker・本番クラウドのいずれでも同じバイナリを
// 使い回せるようにします（getenvを引数で受け取るのは単体テストで
// os.Getenvを本物の環境変数に触れずに差し替えられるようにするためです）。
type config struct {
	Port            string
	DBPath          string
	FrontendOrigin  string
	SecureCookie    bool
	ShutdownTimeout time.Duration
}

func loadConfig(getenv func(string) string) config {
	port := getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/codeforge.db"
	}

	frontendOrigin := getenv("FRONTEND_ORIGIN")
	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:3000"
	}

	// SameSite=None は Secure 必須で、ローカル開発(http)ではCookieが拒否される。
	// 本番(HTTPS)では SECURE_COOKIE=true を設定すること。
	secureCookie := getenv("SECURE_COOKIE") == "true"

	shutdownTimeout := 10 * time.Second
	if v := getenv("SHUTDOWN_TIMEOUT_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			shutdownTimeout = time.Duration(n) * time.Second
		}
	}

	return config{
		Port:            port,
		DBPath:          dbPath,
		FrontendOrigin:  frontendOrigin,
		SecureCookie:    secureCookie,
		ShutdownTimeout: shutdownTimeout,
	}
}

// HealthChecker はヘルスチェックの本体（例: DBへの疎通確認）を表します。
type HealthChecker func() error

func healthHandler(check HealthChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := check(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

// statusRecorder はハンドラが実際に書き込んだステータスコードを記録するための
// http.ResponseWriterラッパーです。アクセスログにステータスを載せるために使います。
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// withRequestLogging は全リクエストの処理結果を構造化ログ（JSON）として出力する。
func withRequestLogging(logger *slog.Logger, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		h.ServeHTTP(rec, r)
		logger.Info("request handled",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

// runHealthCheckClient は `-healthcheck` 引数で起動された場合の挙動です。
// distrolessベースの実行イメージにはシェルもcurlも無いため、Dockerの
// HEALTHCHECK命令から直接このバイナリ自身を呼び出して疎通確認させます。
func runHealthCheckClient() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://localhost:" + port + "/healthz")
	if err != nil {
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-healthcheck" {
		runHealthCheckClient()
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := loadConfig(os.Getenv)

	if err := os.MkdirAll(filepath.Dir(cfg.DBPath), 0o755); err != nil {
		logger.Error("failed to create data directory", "err", err)
		os.Exit(1)
	}

	st, err := store.Open(cfg.DBPath)
	if err != nil {
		logger.Error("failed to open store", "err", err)
		os.Exit(1)
	}
	defer st.Close()

	authSvc := auth.New(st)
	runner := judge.NewRunner()

	// ログイン試行のブルートフォース対策。IPアドレス単位で制限する
	// （単一プロセスでの運用のためインメモリ実装。リバースプロキシ配下での
	// 運用に切り替える場合はクライアントIPの取得方法を見直すこと）。
	loginLimiter := ratelimit.New(10, time.Minute)
	signupLimiter := ratelimit.New(5, time.Minute)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthHandler(st.Ping))
	mux.HandleFunc("GET /api/problems/{id}", getProblemHandler())
	mux.HandleFunc("GET /api/problems/{id}/answer", getAnswerHandler())
	mux.HandleFunc("POST /api/problems/{id}/submit", requireCSRF(submitHandler(runner, authSvc, st)))
	mux.HandleFunc("POST /api/problems/{id}/check", checkHandler())
	mux.HandleFunc("POST /api/problems/{id}/format", formatHandler())
	mux.HandleFunc("GET /api/problems/{id}/draft", getDraftHandler(authSvc, st))
	mux.HandleFunc("POST /api/problems/{id}/draft", requireCSRF(saveDraftHandler(authSvc, st)))

	mux.HandleFunc("POST /api/auth/signup", rateLimited(signupLimiter, requireCSRF(signUpHandler(authSvc, cfg.SecureCookie))))
	mux.HandleFunc("POST /api/auth/login", rateLimited(loginLimiter, requireCSRF(logInHandler(authSvc, cfg.SecureCookie))))
	mux.HandleFunc("POST /api/auth/logout", requireCSRF(logOutHandler(authSvc, cfg.SecureCookie)))
	mux.HandleFunc("GET /api/me", meHandler(authSvc))
	mux.HandleFunc("GET /api/progress", progressHandler(authSvc, st))

	handler := withRequestLogging(logger, withHSTS(cfg.SecureCookie, withCORS(cfg.FrontendOrigin, withCSRFCookie(cfg.SecureCookie, mux))))
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	serveErr := make(chan error, 1)
	go func() {
		logger.Info("listening", "addr", srv.Addr, "frontend_origin", cfg.FrontendOrigin, "db", cfg.DBPath)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serveErr <- err
			return
		}
		serveErr <- nil
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serveErr:
		if err != nil {
			logger.Error("server error", "err", err)
			os.Exit(1)
		}
		return
	case <-ctx.Done():
		stop()
	}

	logger.Info("shutting down", "timeout", cfg.ShutdownTimeout.String())
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "err", err)
	}
}

// withHSTS は Strict-Transport-Security ヘッダーを付与する。
// TLS終端はリバースプロキシ/クラウドLBが担い、このサーバー自身はHTTPで
// 動く構成を前提とするため、r.TLS の有無ではなく SECURE_COOKIE と同じ
// 「本番でHTTPS配信されているか」の判断（環境変数）に乗せる。
func withHSTS(secure bool, h http.Handler) http.Handler {
	if !secure {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		h.ServeHTTP(w, r)
	})
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
			slog.Error("failed to read markdown", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		starter, err := problem.ReadStarter(problemsBaseDir)
		if err != nil {
			slog.Error("failed to read starter code", "err", err)
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
			slog.Error("failed to encode response", "err", err)
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
			slog.Error("failed to read answer code", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(answerResponse{Code: stripBuildIgnoreTag(string(answer))}); err != nil {
			slog.Error("failed to encode response", "err", err)
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
			slog.Error("failed to read go.mod", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		goSum, err := problem.ReadGoSum(problemsBaseDir)
		if err != nil && !os.IsNotExist(err) {
			slog.Error("failed to read go.sum", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		target, err := readTargetIfWritesTest(problem)
		if err != nil {
			slog.Error("failed to read target.go", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), checkTimeout)
		defer cancel()

		diagnostics, err := lint.Check(ctx, goMod, goSum, req.Code, target)
		if err != nil {
			slog.Error("check error", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(checkResponse{Diagnostics: diagnostics}); err != nil {
			slog.Error("failed to encode response", "err", err)
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
			slog.Error("failed to read go.mod", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		goSum, err := problem.ReadGoSum(problemsBaseDir)
		if err != nil && !os.IsNotExist(err) {
			slog.Error("failed to read go.sum", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		target, err := readTargetIfWritesTest(problem)
		if err != nil {
			slog.Error("failed to read target.go", "err", err)
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
			slog.Error("failed to encode response", "err", err)
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
			slog.Error("judge run error", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// ログイン中なら合格を記録する。未ログインでも採点自体は利用できる
		// （その場合、進捗はフロントエンドのlocalStorageにのみ保存される）。
		if result.Passed {
			if user, err := authSvc.UserByToken(sessionToken(r)); err == nil {
				if err := st.MarkProgress(user.ID, problem.ID); err != nil {
					slog.Error("failed to mark progress", "err", err)
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			slog.Error("failed to encode response", "err", err)
		}
	}
}
