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

func (p Problem) readFile(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %q for problem %q: %w", path, p.ID, err)
	}
	return b, nil
}
