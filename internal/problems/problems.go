// Package problems は採点対象の問題定義と、そのファイルへのアクセスを提供します。
package problems

import (
	"fmt"
	"os"
	"path/filepath"
)

// Problem は1つの演習問題を表します。
type Problem struct {
	ID    string
	Title string
	Dir   string // BaseDir からの相対パス
	// WritesTest はユーザーが実装コードではなくテストコードを書く形式の
	// レッスンであることを示す。true の場合、採点は target.go（正しい実装）と
	// mutant.go（わざとバグを仕込んだ実装）の両方に対して行われる
	// （internal/judge の2段階採点を参照）。
	WritesTest bool
}

// allowlist はユーザーから渡される problemID を検証するための許可リストです。
// 任意の文字列をそのままファイルパスに使うとパストラバーサルの危険があるため、
// 必ずこのリストを経由してディレクトリを解決します。
var allowlist = map[string]Problem{
	"web-api-basics-01": {
		ID:    "web-api-basics-01",
		Title: "HTTPサーバーを立てる",
		Dir:   "web-api-basics/01-http-server",
	},
	"web-api-basics-02": {
		ID:    "web-api-basics-02",
		Title: "ルーティングとパスパラメータ",
		Dir:   "web-api-basics/02-routing-and-params",
	},
	"web-api-basics-03": {
		ID:    "web-api-basics-03",
		Title: "リクエストを読み取る",
		Dir:   "web-api-basics/03-reading-requests",
	},
	"web-api-basics-04": {
		ID:    "web-api-basics-04",
		Title: "レスポンスを返す",
		Dir:   "web-api-basics/04-writing-responses",
	},
	"web-api-basics-05": {
		ID:    "web-api-basics-05",
		Title: "共通処理をまとめる",
		Dir:   "web-api-basics/05-middleware",
	},
	"web-api-basics-06": {
		ID:    "web-api-basics-06",
		Title: "道場: TODO管理APIを作ろう",
		Dir:   "web-api-basics/06-capstone-todo-api",
	},
	"concurrency-01": {
		ID:    "concurrency-01",
		Title: "ゴルーチンとWaitGroup",
		Dir:   "concurrency/01-goroutines-and-waitgroup",
	},
	"concurrency-02": {
		ID:    "concurrency-02",
		Title: "channelでやり取りする",
		Dir:   "concurrency/02-channels",
	},
	"concurrency-03": {
		ID:    "concurrency-03",
		Title: "sync.Mutexで排他制御する",
		Dir:   "concurrency/03-mutex",
	},
	"concurrency-04": {
		ID:    "concurrency-04",
		Title: "selectで複数のチャネルを扱う",
		Dir:   "concurrency/04-select",
	},
	"concurrency-05": {
		ID:    "concurrency-05",
		Title: "contextでタイムアウト制御する",
		Dir:   "concurrency/05-context-timeout",
	},
	"concurrency-06": {
		ID:    "concurrency-06",
		Title: "道場: ワーカープールを作ろう",
		Dir:   "concurrency/06-capstone-worker-pool",
	},
	"gin-api-01": {
		ID:    "gin-api-01",
		Title: "Ginの基本",
		Dir:   "gin-api/01-gin-basics",
	},
	"gin-api-02": {
		ID:    "gin-api-02",
		Title: "ルーティンググループとパスパラメータ",
		Dir:   "gin-api/02-routing-groups",
	},
	"gin-api-03": {
		ID:    "gin-api-03",
		Title: "リクエストのバインディングとバリデーション",
		Dir:   "gin-api/03-binding-validation",
	},
	"gin-api-04": {
		ID:    "gin-api-04",
		Title: "Ginのミドルウェア",
		Dir:   "gin-api/04-middleware",
	},
	"gin-api-05": {
		ID:    "gin-api-05",
		Title: "道場: GinでCRUD APIを作ろう",
		Dir:   "gin-api/05-capstone-crud",
	},
	"database-01": {
		ID:    "database-01",
		Title: "データベースに接続してテーブルを作る",
		Dir:   "database/01-connect-and-create-table",
	},
	"database-02": {
		ID:    "database-02",
		Title: "データを登録・取得する",
		Dir:   "database/02-insert-and-select",
	},
	"database-03": {
		ID:    "database-03",
		Title: "データを更新・削除する",
		Dir:   "database/03-update-and-delete",
	},
	"database-04": {
		ID:    "database-04",
		Title: "GORMの基本",
		Dir:   "database/04-gorm-basics",
	},
	"database-05": {
		ID:    "database-05",
		Title: "道場: GORMでCRUDを完成させよう",
		Dir:   "database/05-capstone-gorm-crud",
	},
	"auth-jwt-01": {
		ID:    "auth-jwt-01",
		Title: "パスワードをハッシュ化する",
		Dir:   "auth-jwt/01-password-hashing",
	},
	"auth-jwt-02": {
		ID:    "auth-jwt-02",
		Title: "JWTを発行する",
		Dir:   "auth-jwt/02-issue-jwt",
	},
	"auth-jwt-03": {
		ID:    "auth-jwt-03",
		Title: "JWTを検証する",
		Dir:   "auth-jwt/03-verify-jwt",
	},
	"auth-jwt-04": {
		ID:    "auth-jwt-04",
		Title: "認証ミドルウェアでAPIを保護する",
		Dir:   "auth-jwt/04-auth-middleware",
	},
	"auth-jwt-05": {
		ID:    "auth-jwt-05",
		Title: "道場: ログインAPIを作ろう",
		Dir:   "auth-jwt/05-capstone-login-api",
	},
	"testing-01": {
		ID:         "testing-01",
		Title:      "テーブル駆動テストを書く",
		Dir:        "testing/01-table-driven-tests",
		WritesTest: true,
	},
	"testing-02": {
		ID:         "testing-02",
		Title:      "サブテストを使いこなす",
		Dir:        "testing/02-subtests",
		WritesTest: true,
	},
	"testing-03": {
		ID:         "testing-03",
		Title:      "インターフェースでモックする",
		Dir:        "testing/03-interfaces-and-mocks",
		WritesTest: true,
	},
	"testing-04": {
		ID:         "testing-04",
		Title:      "エラーケースをテストする",
		Dir:        "testing/04-testing-errors",
		WritesTest: true,
	},
	"testing-05": {
		ID:         "testing-05",
		Title:      "道場: HTTPハンドラーをテストする",
		Dir:        "testing/05-capstone-http-handler",
		WritesTest: true,
	},
	"error-handling-01": {
		ID:    "error-handling-01",
		Title: "カスタムエラー型を定義する",
		Dir:   "error-handling/01-custom-error-types",
	},
	"error-handling-02": {
		ID:    "error-handling-02",
		Title: "エラーをラップする",
		Dir:   "error-handling/02-wrapping-errors",
	},
	"error-handling-03": {
		ID:    "error-handling-03",
		Title: "errors.Is / errors.As で判定する",
		Dir:   "error-handling/03-errors-is-and-as",
	},
	"error-handling-04": {
		ID:    "error-handling-04",
		Title: "複数のエラーをまとめる",
		Dir:   "error-handling/04-joining-errors",
	},
	"error-handling-05": {
		ID:    "error-handling-05",
		Title: "道場: 実践的なエラーハンドリング",
		Dir:   "error-handling/05-capstone-order-processing",
	},
	"architecture-01": {
		ID:    "architecture-01",
		Title: "インターフェースで依存を注入する",
		Dir:   "architecture/01-dependency-injection",
	},
	"architecture-02": {
		ID:    "architecture-02",
		Title: "リポジトリパターンでデータアクセスを抽象化する",
		Dir:   "architecture/02-repository-pattern",
	},
	"architecture-03": {
		ID:    "architecture-03",
		Title: "サービス層でビジネスロジックを分離する",
		Dir:   "architecture/03-service-layer",
	},
	"architecture-04": {
		ID:    "architecture-04",
		Title: "関数オプションパターンで柔軟な初期化をする",
		Dir:   "architecture/04-functional-options",
	},
	"architecture-05": {
		ID:    "architecture-05",
		Title: "道場: レイヤードアーキテクチャを組み立てる",
		Dir:   "architecture/05-capstone-layered-app",
	},
	"deploy-01": {
		ID:    "deploy-01",
		Title: "環境変数で設定を切り替える",
		Dir:   "deploy/01-env-config",
	},
	"deploy-02": {
		ID:    "deploy-02",
		Title: "ヘルスチェックエンドポイントを実装する",
		Dir:   "deploy/02-health-check",
	},
	"deploy-03": {
		ID:    "deploy-03",
		Title: "グレースフルシャットダウンを実装する",
		Dir:   "deploy/03-graceful-shutdown",
	},
	"deploy-04": {
		ID:    "deploy-04",
		Title: "構造化ロギングでコンテナ環境に対応する",
		Dir:   "deploy/04-structured-logging",
	},
	"deploy-05": {
		ID:    "deploy-05",
		Title: "道場: 本番向けのHTTPサーバーを組み立てる",
		Dir:   "deploy/05-capstone-production-server",
	},
}

// Get は problemID に対応する Problem を返します。存在しない場合は ok が false になります。
func Get(id string) (Problem, bool) {
	p, ok := allowlist[id]
	return p, ok
}

// TestFilePath は非公開の採点用テストファイルの絶対パスを返します。
func (p Problem) TestFilePath(baseDir string) string {
	return filepath.Join(baseDir, p.Dir, "solution_test.go")
}

// MarkdownPath は問題説明用Markdownファイルのパスを返します。
func (p Problem) MarkdownPath(baseDir string) string {
	return filepath.Join(baseDir, p.Dir, "problem.md")
}

// StarterPath はユーザーに提示する雛形コードのパスを返します。
func (p Problem) StarterPath(baseDir string) string {
	return filepath.Join(baseDir, p.Dir, "starter.go")
}

// AnswerPath は模範解答コードのパスを返します。
func (p Problem) AnswerPath(baseDir string) string {
	return filepath.Join(baseDir, p.Dir, "answer.go")
}

// GoModPath は採点ワークスペースに配置する go.mod のパスを返します。
// 標準ライブラリのみで完結する問題は `module submission` / `go 1.22` だけの
// 最小限の内容、外部モジュールに依存する問題は go mod tidy で生成した内容を置く。
func (p Problem) GoModPath(baseDir string) string {
	return filepath.Join(baseDir, p.Dir, "go.mod")
}

// GoSumPath は採点ワークスペースに配置する go.sum のパスを返します。
// 標準ライブラリのみで完結する問題には存在しません（ReadGoSum が os.ErrNotExist を返します）。
func (p Problem) GoSumPath(baseDir string) string {
	return filepath.Join(baseDir, p.Dir, "go.sum")
}

// TargetPath は（WritesTest問題の）正しい実装のパスを返します。
func (p Problem) TargetPath(baseDir string) string {
	return filepath.Join(baseDir, p.Dir, "target.go")
}

// MutantPath は（WritesTest問題の）わざとバグを仕込んだ実装のパスを返します。
func (p Problem) MutantPath(baseDir string) string {
	return filepath.Join(baseDir, p.Dir, "mutant.go")
}

// ReadTestFile は採点用テストファイルの中身を読み込みます。
func (p Problem) ReadTestFile(baseDir string) ([]byte, error) {
	return p.readFile(p.TestFilePath(baseDir))
}

// ReadMarkdown は問題説明用Markdownの中身を読み込みます。
func (p Problem) ReadMarkdown(baseDir string) ([]byte, error) {
	return p.readFile(p.MarkdownPath(baseDir))
}

// ReadStarter はユーザーに提示する雛形コードの中身を読み込みます。
func (p Problem) ReadStarter(baseDir string) ([]byte, error) {
	return p.readFile(p.StarterPath(baseDir))
}

// ReadAnswer は模範解答コードの中身を読み込みます。
func (p Problem) ReadAnswer(baseDir string) ([]byte, error) {
	return p.readFile(p.AnswerPath(baseDir))
}

// ReadGoMod は採点ワークスペース用 go.mod の中身を読み込みます。
func (p Problem) ReadGoMod(baseDir string) ([]byte, error) {
	return p.readFile(p.GoModPath(baseDir))
}

// ReadGoSum は採点ワークスペース用 go.sum の中身を読み込みます。
// ファイルが存在しない場合（標準ライブラリのみの問題）は os.IsNotExist で判定できるエラーを返します。
func (p Problem) ReadGoSum(baseDir string) ([]byte, error) {
	b, err := os.ReadFile(p.GoSumPath(baseDir))
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ReadTarget は（WritesTest問題の）正しい実装の中身を読み込みます。
func (p Problem) ReadTarget(baseDir string) ([]byte, error) {
	return p.readFile(p.TargetPath(baseDir))
}

// ReadMutant は（WritesTest問題の）わざとバグを仕込んだ実装の中身を読み込みます。
func (p Problem) ReadMutant(baseDir string) ([]byte, error) {
	return p.readFile(p.MutantPath(baseDir))
}

func (p Problem) readFile(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %q for problem %q: %w", path, p.ID, err)
	}
	return b, nil
}
