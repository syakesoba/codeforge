// Package lint はエディタ向けの軽量な静的チェック（コンパイルエラー検知・
// import自動修正）を提供します。
//
// go build はコンパイルのみでユーザーコードを一切実行しない
// （go run/go testと違い init() や main() が走らない）ため、
// 採点用のDockerサンドボックス（internal/judge）とは異なり、
// バックエンドプロセスから直接実行しても安全です。
// これによりDockerの起動オーバーヘッドが無くなり、
// エディタ入力に対するリアルタイムに近い検知が可能になります。
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

// diagnosticPattern は `go build` が標準エラー出力に書く1行分の形式
// （例: "./solution.go:5:2: undefined: fmt"）にマッチします。
var diagnosticPattern = regexp.MustCompile(`^(?:\./)?([^:]+):(\d+):(\d+):\s*(.+)$`)

const checkTimeout = 8 * time.Second

// setupWorkspace はレッスンのgo.mod/go.sumとユーザーコードを一時ディレクトリに
// 配置し、そのディレクトリパスを返します。呼び出し側で削除してください。
//
// go.mod を実ファイルとして置くのは、goimportsの依存解決（Gin/GORM等の
// パッケージ名からimportパスを引く処理）が「対象ファイルのディレクトリから
// 一番近いgo.mod」を辿ってモジュールの依存グラフを見に行くためで、
// ディレクトリと無関係なパスを渡すと解決に失敗する。
func setupWorkspace(goMod, goSum []byte, code string) (dir string, err error) {
	workdir, err := os.MkdirTemp("", "lint-*")
	if err != nil {
		return "", fmt.Errorf("failed to create workspace: %w", err)
	}

	files := map[string][]byte{
		"go.mod":      goMod,
		"solution.go": []byte(code),
	}
	if len(goSum) > 0 {
		files["go.sum"] = goSum
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
// `go build` のみを実行してコンパイルエラー（未インポート・未使用変数等）を収集します。
// 採点用の非公開テストは含めません。
func Check(ctx context.Context, goMod, goSum []byte, code string) ([]Diagnostic, error) {
	workdir, err := setupWorkspace(goMod, goSum, code)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(workdir)

	return buildDiagnostics(ctx, workdir, []byte(code))
}

// buildDiagnostics は既存のワークスペース（go.mod/go.sum配置済み）に
// solution.go を書き直して go build し、診断を返す。
// 同時実行数はセマフォで制限する。
func buildDiagnostics(ctx context.Context, workdir string, code []byte) ([]Diagnostic, error) {
	select {
	case checkSemaphore <- struct{}{}:
		defer func() { <-checkSemaphore }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if err := os.WriteFile(filepath.Join(workdir, "solution.go"), code, 0o644); err != nil {
		return nil, fmt.Errorf("failed to write solution.go: %w", err)
	}

	buildCtx, cancel := context.WithTimeout(ctx, checkTimeout)
	defer cancel()

	cmd := exec.CommandContext(buildCtx, "go", "build", "./...")
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
