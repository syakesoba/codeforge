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
	// WorkDir は採点ワークスペースを作る親ディレクトリ。空ならOSの一時ディレクトリを使う。
	//
	// バックエンド自身をコンテナで動かし、ホストのDockerデーモン（ソケットをマウント）に
	// 採点コンテナを起動させる構成では、`docker run -v` に渡すパスはホスト側で解決される。
	// そのため、ホストとバックエンドのコンテナで同じパスになるようバインドマウントした
	// ディレクトリを指定する必要がある。
	WorkDir string
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

// runOutcome は runOnce 1回分（docker run 1回分）の実行結果です。
type runOutcome struct {
	passed   bool
	output   string
	duration time.Duration
	timedOut bool
	// buildFailed はテストを1件も実行する前にコンパイルが失敗したことを表す。
	buildFailed bool
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
	supportFiles, err := problem.ReadSupportFiles(r.ProblemsBaseDir)
	if err != nil {
		return Result{}, err
	}

	if problem.WritesTest {
		return r.runWritesTest(ctx, problem, code, goMod, goSum, supportFiles)
	}
	return r.runImplementation(ctx, problem, code, goMod, goSum, supportFiles)
}

// runImplementation は通常のレッスン（ユーザーが実装コードを書き、非公開テストで
// 検証する形式）を採点する。
func (r Runner) runImplementation(ctx context.Context, problem problems.Problem, code string, goMod, goSum []byte, supportFiles map[string]string) (Result, error) {
	testSrc, err := problem.ReadTestFile(r.ProblemsBaseDir)
	if err != nil {
		return Result{}, err
	}

	files := cloneFiles(supportFiles)
	files["solution.go"] = code
	files["solution_test.go"] = string(testSrc)

	out, err := r.runOnce(ctx, goMod, goSum, files)
	if err != nil {
		return Result{}, err
	}
	if out.timedOut {
		return Result{
			Passed:     false,
			Output:     "実行時間の上限を超えたため停止しました。",
			DurationMs: out.duration.Milliseconds(),
			Error:      "timeout",
		}, nil
	}

	return Result{Passed: out.passed, Output: out.output, DurationMs: out.duration.Milliseconds()}, nil
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
func (r Runner) runWritesTest(ctx context.Context, problem problems.Problem, code string, goMod, goSum []byte, supportFiles map[string]string) (Result, error) {
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

	start := time.Now()

	stage1, err := r.runOnce(ctx, goMod, goSum, writesTestStageFiles(supportFiles, code, string(hiddenTest), string(target)))
	if err != nil {
		return Result{}, err
	}
	if stage1.timedOut {
		return Result{
			Passed:     false,
			Output:     "実行時間の上限を超えたため停止しました。",
			DurationMs: time.Since(start).Milliseconds(),
			Error:      "timeout",
		}, nil
	}
	if !stage1.passed {
		return Result{
			Passed:     false,
			Output:     stage1.output,
			DurationMs: time.Since(start).Milliseconds(),
		}, nil
	}

	stage2, err := r.runOnce(ctx, goMod, goSum, writesTestStageFiles(supportFiles, code, string(hiddenTest), string(mutant)))
	if err != nil {
		return Result{}, err
	}
	duration := time.Since(start)
	if stage2.timedOut {
		return Result{
			Passed:     false,
			Output:     "実行時間の上限を超えたため停止しました。",
			DurationMs: duration.Milliseconds(),
			Error:      "timeout",
		}, nil
	}
	if stage2.buildFailed {
		// Stage 1 と同じユーザーコード・非公開テストでコンパイルできているため、
		// Stage 2 だけビルドに失敗するのは mutant.go 側（サーバーのコンテンツ）の不備。
		// これを「バグを検出できた」と扱うと、どんなテストでも合格してしまうため
		// サーバーエラーとして返す。
		return Result{}, fmt.Errorf("stage 2 (mutant) failed to build for problem %q; check mutant.go:\n%s", problem.ID, stage2.output)
	}
	if stage2.passed {
		return Result{
			Passed: false,
			Output: "テストは実行できましたが、わざとバグを仕込んだ実装でも合格してしまいました。\n" +
				"期待値との比較や、一致しない場合に t.Errorf / t.Fatalf で失敗させる処理が" +
				"正しく書けているか確認してください。\n\n" + stage2.output,
			DurationMs: duration.Milliseconds(),
			Error:      "insufficient_test",
		}, nil
	}

	return Result{Passed: true, Output: stage1.output, DurationMs: duration.Milliseconds()}, nil
}

// writesTestStageFiles は WritesTest レッスンの1ステージ分のワークスペースに置く
// ファイル（go.mod/go.sum を除く）を組み立てる。impl には Stage 1 では正しい実装、
// Stage 2 ではバグを仕込んだ実装を渡し、どちらも target.go として配置する。
func writesTestStageFiles(supportFiles map[string]string, code, hiddenTest, impl string) map[string]string {
	files := cloneFiles(supportFiles)
	files["solution_test.go"] = code
	files["hidden_test.go"] = hiddenTest
	files["target.go"] = impl
	return files
}

// runOnce は1回分のワークスペースを作成し、Dockerサンドボックス内で
// `go test -vet=off -json ./...` を実行する。
func (r Runner) runOnce(ctx context.Context, goMod, goSum []byte, files map[string]string) (runOutcome, error) {
	workdir, err := os.MkdirTemp(r.WorkDir, "judge-*")
	if err != nil {
		return runOutcome{}, fmt.Errorf("failed to create workspace: %w", err)
	}
	defer os.RemoveAll(workdir)

	// MkdirTemp は所有者だけが読める 0700 で作る。採点コンテナは別ユーザー
	// （採点イメージの judge, uid 1000）で動くため、読めるように広げておく。
	if err := os.Chmod(workdir, 0o755); err != nil {
		return runOutcome{}, fmt.Errorf("failed to chmod workspace: %w", err)
	}

	if err := writeWorkspace(workdir, goMod, goSum, files); err != nil {
		return runOutcome{}, err
	}

	runner := dockerRunner{Image: r.Image}
	start := time.Now()
	stdout, stderr, runErr := runner.run(ctx, workdir, r.Timeout)
	duration := time.Since(start)

	if errors.Is(runErr, errTimeout) {
		return runOutcome{duration: duration, timedOut: true}, nil
	}

	// go test は「全テストがpassしたときだけ終了コード0」という契約を守るため、
	// 合否判定はプロセスの終了コード（runErr）を正とする。JSON解析は表示用ログの
	// 整形にのみ使う（fail イベントが出る前にプロセスがpanicで落ちるケースなどを
	// 正しく不合格として扱うため）。
	return runOutcome{
		passed:      runErr == nil,
		output:      formatGoTestOutput(stdout, stderr),
		duration:    duration,
		buildFailed: runErr != nil && isBuildFailure(stdout, stderr),
	}, nil
}

// writeWorkspace は go.mod・go.sum（空でなければ）・files を dir に書き込む。
func writeWorkspace(dir string, goMod, goSum []byte, files map[string]string) error {
	allFiles := map[string]string{"go.mod": string(goMod)}
	if len(goSum) > 0 {
		allFiles["go.sum"] = string(goSum)
	}
	for name, content := range files {
		allFiles[name] = content
	}

	for name, content := range allFiles {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			return fmt.Errorf("failed to write %s: %w", name, err)
		}
	}
	return nil
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
