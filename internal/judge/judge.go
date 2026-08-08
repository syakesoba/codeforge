// Package judge はユーザーが投稿したGoコードを、Dockerサンドボックス内で
// 問題ごとの非公開テストとともに実行し、合否を判定します。
package judge

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/syakesoba/codeforge/internal/problems"
)

// Result は採点結果です。
type Result struct {
	Passed     bool   `json:"passed"`
	Output     string `json:"output"`
	DurationMs int64  `json:"durationMs"`
	Error      string `json:"error,omitempty"`
}

// Runner はサンドボックス実行の設定を保持します。
type Runner struct {
	ProblemsBaseDir string
	Image           string
	Timeout         time.Duration
}

// NewRunner はデフォルト設定の Runner を作成します。
func NewRunner() Runner {
	return Runner{
		ProblemsBaseDir: "problems",
		Image:           "codeforge-judge:latest",
		// 外部モジュール（Gin/GORM等）を使うレッスンはリンクに時間がかかるため、
		// stdlibのみの場合より余裕を持たせている。
		Timeout: 30 * time.Second,
	}
}

// Run はユーザーコードを problem の非公開テストとともにサンドボックス実行します。
func (r Runner) Run(ctx context.Context, problem problems.Problem, code string) (Result, error) {
	goMod, err := problem.ReadGoMod(r.ProblemsBaseDir)
	if err != nil {
		return Result{}, err
	}

	testSrc, err := problem.ReadTestFile(r.ProblemsBaseDir)
	if err != nil {
		return Result{}, err
	}

	workdir, err := os.MkdirTemp("", "judge-*")
	if err != nil {
		return Result{}, fmt.Errorf("failed to create workspace: %w", err)
	}
	defer os.RemoveAll(workdir)

	files := map[string]string{
		"go.mod":           string(goMod),
		"solution.go":      code,
		"solution_test.go": string(testSrc),
	}

	// go.sum は外部モジュールに依存する問題にのみ存在する（stdlibのみの問題では不要）
	if goSum, err := problem.ReadGoSum(r.ProblemsBaseDir); err == nil {
		files["go.sum"] = string(goSum)
	} else if !os.IsNotExist(err) {
		return Result{}, fmt.Errorf("failed to read go.sum: %w", err)
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(workdir, name), []byte(content), 0o644); err != nil {
			return Result{}, fmt.Errorf("failed to write %s: %w", name, err)
		}
	}

	runner := dockerRunner{Image: r.Image}
	start := time.Now()
	stdout, stderr, runErr := runner.run(ctx, workdir, r.Timeout)
	duration := time.Since(start)

	if errors.Is(runErr, errTimeout) {
		return Result{
			Passed:     false,
			Output:     "実行時間の上限を超えたため停止しました。",
			DurationMs: duration.Milliseconds(),
			Error:      "timeout",
		}, nil
	}

	// go test は「全テストがpassしたときだけ終了コード0」という契約を守るため、
	// 合否判定はプロセスの終了コード（runErr）を正とする。JSON解析は表示用ログの
	// 整形にのみ使う（fail イベントが出る前にプロセスがpanicで落ちるケースなどを
	// 正しく不合格として扱うため）。
	output := formatGoTestOutput(stdout, stderr)
	return Result{
		Passed:     runErr == nil,
		Output:     output,
		DurationMs: duration.Milliseconds(),
	}, nil
}
