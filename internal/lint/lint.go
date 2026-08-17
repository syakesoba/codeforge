// Package lint はエディタ向けの軽量な静的チェック（コンパイルエラー検知・
// import自動修正）を提供します。
//
// go build はコンパイルのみでユーザーコードを一切実行しない
// （go run/go testと違い init() や main() が走らない）ため、
// 採点用のDockerサンドボックス（internal/judge）とは異なり、
// バックエンドプロセスから直接実行しても安全です。
// これによりDockerの起動オーバーヘッドが無くなり、
// エディタ入力に対するリアルタイムに近い検知が可能になります。
//
// 「テストを書く」形式のレッスン（internal/problems の WritesTest）では、
// ユーザーコードは実装ではなくテストコードなので、`go build` の対象にならない
// （_test.goファイルはgo buildの対象外）。この場合は target（正しい実装）を
// 同じワークスペースに置いた上で、ユーザーコードを solution_test.go として配置し、
// `go test -run=^$` で（何も実行せず）コンパイルのみ行う。
package lint

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// maxConcurrentChecks は同時に実行する go build プロセス数の上限です。
// 大量の同時チェックでホストのCPU/メモリが枯渇しないよう制限します。
const maxConcurrentChecks = 4

var checkSemaphore = make(chan struct{}, maxConcurrentChecks)

// Diagnostic はエディタに表示する1件の診断情報です。
type Diagnostic struct {
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Message string `json:"message"`
}

// diagnosticPattern は `go build`/`go test` が標準エラー出力に書く1行分の形式
// （例: "./solution.go:5:2: undefined: fmt"）にマッチします。
var diagnosticPattern = regexp.MustCompile(`^(?:\./)?([^:]+):(\d+):(\d+):\s*(.+)$`)

const checkTimeout = 8 * time.Second

// codeFilename は target の有無に応じて、ユーザーコードを配置するファイル名を返す。
// target が指定されている場合（WritesTestレッスン）はテストとして実行されるよう
// "_test.go" にする。
func codeFilename(target []byte) string {
	if target != nil {
		return "solution_test.go"
	}
	return "solution.go"
}

// setupWorkspace はレッスンのgo.mod/go.sumとユーザーコードを一時ディレクトリに
// 配置し、そのディレクトリパスを返します。呼び出し側で削除してください。
//
// go.mod を実ファイルとして置くのは、goimportsの依存解決（Gin/GORM等の
// パッケージ名からimportパスを引く処理）が「対象ファイルのディレクトリから
// 一番近いgo.mod」を辿ってモジュールの依存グラフを見に行くためで、
// ディレクトリと無関係なパスを渡すと解決に失敗する。
//
// target が非nilの場合（WritesTestレッスン）は、ユーザーコードが参照する
// 「正しい実装」を target.go として同じワークスペースに配置する。
func setupWorkspace(goMod, goSum []byte, code string, target []byte) (dir string, err error) {
	workdir, err := os.MkdirTemp("", "lint-*")
	if err != nil {
		return "", fmt.Errorf("failed to create workspace: %w", err)
	}

	files := map[string][]byte{
		"go.mod":             goMod,
		codeFilename(target): []byte(code),
	}
	if len(goSum) > 0 {
		files["go.sum"] = goSum
	}
	if target != nil {
		files["target.go"] = target
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(workdir, name), content, 0o644); err != nil {
			os.RemoveAll(workdir)
			return "", fmt.Errorf("failed to write %s: %w", name, err)
		}
	}
	return workdir, nil
}

// Check はユーザーコードを go.mod/go.sum 付きの一時ワークスペースに配置し、
// `go build`（WritesTestレッスンでは `go test -run=^$`）のみを実行して
// コンパイルエラー（未インポート・未使用変数等）を収集します。
// 採点用の非公開テストは含めません。
func Check(ctx context.Context, goMod, goSum []byte, code string, target []byte) ([]Diagnostic, error) {
	workdir, err := setupWorkspace(goMod, goSum, code, target)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(workdir)

	return buildDiagnostics(ctx, workdir, []byte(code), target != nil)
}

// buildDiagnostics は既存のワークスペース（go.mod/go.sum配置済み）にユーザーコードを
// 書き直してコンパイルし、診断を返す。同時実行数はセマフォで制限する。
func buildDiagnostics(ctx context.Context, workdir string, code []byte, writesTest bool) ([]Diagnostic, error) {
	select {
	case checkSemaphore <- struct{}{}:
		defer func() { <-checkSemaphore }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	filename := "solution.go"
	args := []string{"build", "./..."}
	if writesTest {
		filename = "solution_test.go"
		// -run=^$ はどのテストにもマッチしないため、コンパイルのみ行い何も実行しない。
		args = []string{"test", "-run=^$", "-vet=off", "./..."}
	}

	if err := os.WriteFile(filepath.Join(workdir, filename), code, 0o644); err != nil {
		return nil, fmt.Errorf("failed to write %s: %w", filename, err)
	}

	buildCtx, cancel := context.WithTimeout(ctx, checkTimeout)
	defer cancel()

	cmd := exec.CommandContext(buildCtx, "go", args...)
	cmd.Dir = workdir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	_ = cmd.Run() // 非0終了はコンパイルエラーとして扱うため無視する

	return parseDiagnostics(stderr.String()), nil
}

func parseDiagnostics(stderr string) []Diagnostic {
	diagnostics := []Diagnostic{}
	for _, line := range strings.Split(stderr, "\n") {
		m := diagnosticPattern.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		lineNum, err := strconv.Atoi(m[2])
		if err != nil {
			continue
		}
		col, err := strconv.Atoi(m[3])
		if err != nil {
			continue
		}
		diagnostics = append(diagnostics, Diagnostic{
			Line:    lineNum,
			Column:  col,
			Message: m[4],
		})
	}
	return diagnostics
}
