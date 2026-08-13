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

// maxConcurrentSubmissions は同時に起動するDockerコンテナ数の上限です。
// 1コンテナあたり --cpus 1.0 / --memory 768m を使うため、ホストのCPU/メモリを
// 考慮して同時実行数を絞る（想定: 4コア程度の小規模ホストでバックエンド・
// フロントエンド分の余力を残す）。上限を超えた分はチャネルの空きを待つ
// （＝キューイングされる）。
const maxConcurrentSubmissions = 3

var submitSemaphore = make(chan struct{}, maxConcurrentSubmissions)

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
// 同時実行数が上限に達している間はキューで待機し、待機中にctxがタイムアウト
// した場合はエラーではなく不合格のResultとして返す（フロントエンドの表示を
// 通常のタイムアウトと同様に扱えるようにするため）。
func (r Runner) Run(ctx context.Context, problem problems.Problem, code string) (Result, error) {
	waitStart := time.Now()
	select {
	case submitSemaphore <- struct{}{}:
		defer func() { <-submitSemaphore }()
	case <-ctx.Done():
		return Result{
			Passed:     false,
			Output:     "採点の混雑によりキューで待機中にタイムアウトしました。しばらくしてからもう一度お試しください。",
			DurationMs: time.Since(waitStart).Milliseconds(),
			Error:      "queue_timeout",
		}, nil
	}

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
