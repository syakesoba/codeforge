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
	goSum, err := readGoSumOrNil(problem, r.ProblemsBaseDir)
	if err != nil {
		return Result{}, err
	}

	if problem.WritesTest {
		return r.runWritesTest(ctx, problem, code, goMod, goSum)
	}
	return r.runImplementation(ctx, problem, code, goMod, goSum)
}

// runImplementation は通常のレッスン（ユーザーが実装コードを書き、非公開テストで
// 検証する形式）を採点する。
func (r Runner) runImplementation(ctx context.Context, problem problems.Problem, code string, goMod, goSum []byte) (Result, error) {
	testSrc, err := problem.ReadTestFile(r.ProblemsBaseDir)
	if err != nil {
		return Result{}, err
	}

	files := map[string]string{
		"solution.go":      code,
		"solution_test.go": string(testSrc),
	}

	passed, output, duration, timedOut, err := r.runOnce(ctx, goMod, goSum, files)
	if err != nil {
		return Result{}, err
	}
	if timedOut {
		return Result{
			Passed:     false,
			Output:     "実行時間の上限を超えたため停止しました。",
			DurationMs: duration.Milliseconds(),
			Error:      "timeout",
		}, nil
	}

	return Result{Passed: passed, Output: output, DurationMs: duration.Milliseconds()}, nil
}

// runWritesTest は「テストを書く」レッスン（ユーザーがテストコードを書き、
// 与えられた実装を検証する形式）を採点する。
//
// ユーザーのテストコードは、比較や失敗処理を一切書かなくても構文上は
// 合格してしまう（何もしない関数でも「非公開テストが失敗しなかった」ことに
// なってしまう）。これを防ぐため、2段階で採点する:
//
//  1. 正しい実装（target.go）に対して実行し、合格することを確認する
//  2. わざとバグを仕込んだ実装（mutant.go）に対して実行し、
//     今度は不合格になる（＝ユーザーのテストがバグを検出できる）ことを確認する
//
// 両方を満たして初めて合格とする。
func (r Runner) runWritesTest(ctx context.Context, problem problems.Problem, code string, goMod, goSum []byte) (Result, error) {
	target, err := problem.ReadTarget(r.ProblemsBaseDir)
	if err != nil {
		return Result{}, err
	}
	mutant, err := problem.ReadMutant(r.ProblemsBaseDir)
	if err != nil {
		return Result{}, err
	}
	hiddenTest, err := problem.ReadTestFile(r.ProblemsBaseDir)
	if err != nil {
		return Result{}, err
	}

	baseFiles := map[string]string{
		"solution_test.go": code,
		"hidden_test.go":   string(hiddenTest),
	}

	start := time.Now()

	stage1Files := cloneFiles(baseFiles)
	stage1Files["target.go"] = string(target)
	stage1Passed, stage1Output, _, stage1TimedOut, err := r.runOnce(ctx, goMod, goSum, stage1Files)
	if err != nil {
		return Result{}, err
	}
	if stage1TimedOut {
		return Result{
			Passed:     false,
			Output:     "実行時間の上限を超えたため停止しました。",
			DurationMs: time.Since(start).Milliseconds(),
			Error:      "timeout",
		}, nil
	}
	if !stage1Passed {
		return Result{
			Passed:     false,
			Output:     stage1Output,
			DurationMs: time.Since(start).Milliseconds(),
		}, nil
	}

	stage2Files := cloneFiles(baseFiles)
	stage2Files["target.go"] = string(mutant)
	stage2Passed, stage2Output, _, stage2TimedOut, err := r.runOnce(ctx, goMod, goSum, stage2Files)
	if err != nil {
		return Result{}, err
	}
	duration := time.Since(start)
	if stage2TimedOut {
		return Result{
			Passed:     false,
			Output:     "実行時間の上限を超えたため停止しました。",
			DurationMs: duration.Milliseconds(),
			Error:      "timeout",
		}, nil
	}
	if stage2Passed {
		return Result{
			Passed: false,
			Output: "テストは実行できましたが、わざとバグを仕込んだ実装でも合格してしまいました。\n" +
				"期待値との比較や、一致しない場合に t.Errorf / t.Fatalf で失敗させる処理が" +
				"正しく書けているか確認してください。\n\n" + stage2Output,
			DurationMs: duration.Milliseconds(),
			Error:      "insufficient_test",
		}, nil
	}

	return Result{Passed: true, Output: stage1Output, DurationMs: duration.Milliseconds()}, nil
}

// runOnce は1回分のワークスペースを作成し、Dockerサンドボックス内で
// `go test -vet=off -json ./...` を実行する。
func (r Runner) runOnce(ctx context.Context, goMod, goSum []byte, files map[string]string) (passed bool, output string, duration time.Duration, timedOut bool, err error) {
	workdir, err := os.MkdirTemp("", "judge-*")
	if err != nil {
		return false, "", 0, false, fmt.Errorf("failed to create workspace: %w", err)
	}
	defer os.RemoveAll(workdir)

	allFiles := map[string]string{"go.mod": string(goMod)}
	if len(goSum) > 0 {
		allFiles["go.sum"] = string(goSum)
	}
	for name, content := range files {
		allFiles[name] = content
	}

	for name, content := range allFiles {
		if err := os.WriteFile(filepath.Join(workdir, name), []byte(content), 0o644); err != nil {
			return false, "", 0, false, fmt.Errorf("failed to write %s: %w", name, err)
		}
	}

	runner := dockerRunner{Image: r.Image}
	start := time.Now()
	stdout, stderr, runErr := runner.run(ctx, workdir, r.Timeout)
	duration = time.Since(start)

	if errors.Is(runErr, errTimeout) {
		return false, "", duration, true, nil
	}

	// go test は「全テストがpassしたときだけ終了コード0」という契約を守るため、
	// 合否判定はプロセスの終了コード（runErr）を正とする。JSON解析は表示用ログの
	// 整形にのみ使う（fail イベントが出る前にプロセスがpanicで落ちるケースなどを
	// 正しく不合格として扱うため）。
	output = formatGoTestOutput(stdout, stderr)
	return runErr == nil, output, duration, false, nil
}

func readGoSumOrNil(problem problems.Problem, baseDir string) ([]byte, error) {
	goSum, err := problem.ReadGoSum(baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read go.sum: %w", err)
	}
	return goSum, nil
}

func cloneFiles(m map[string]string) map[string]string {
	c := make(map[string]string, len(m))
	for k, v := range m {
		c[k] = v
	}
	return c
}
