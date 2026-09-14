package judge

import (
	"bytes"
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/syakesoba/codeforge/internal/problems"
)

// allWritesTestIDs は許可リストに登録されている「テストを書く」レッスンのID。
var allWritesTestIDs = []string{"testing-01", "testing-02", "testing-03", "testing-04", "testing-05"}

// TestWritesTestContent は、各 WritesTest レッスンのコンテンツ（target.go /
// mutant.go / 非公開テスト / 模範解答）が2段階採点の前提を満たしているかを、
// Docker を使わずホストの go コマンドで確認する。
//
//   - Stage 1: 模範解答のテストは正しい実装で合格する
//   - Stage 2: 模範解答のテストはバグ入り実装で「テスト失敗」になる
//     （ビルドエラーで落ちるのはコンテンツの不備。以前は mutant.go の
//     `//go:build ignore` が残っていて、どんなテストでも合格していた）
func TestWritesTestContent(t *testing.T) {
	if testing.Short() {
		t.Skip("go test を実際に実行するため -short ではスキップする")
	}
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go コマンドが見つからない")
	}

	for _, id := range allWritesTestIDs {
		t.Run(id, func(t *testing.T) {
			problem, ok := problems.Get(id)
			if !ok || !problem.WritesTest {
				t.Fatalf("%s は WritesTest レッスンとして登録されていない", id)
			}

			goMod, err := problem.ReadGoMod(problemsDirForTest)
			if err != nil {
				t.Fatal(err)
			}
			goSum, err := readGoSumOrNil(problem, problemsDirForTest)
			if err != nil {
				t.Fatal(err)
			}
			support, err := problem.ReadSupportFiles(problemsDirForTest)
			if err != nil {
				t.Fatal(err)
			}
			target, err := problem.ReadTarget(problemsDirForTest)
			if err != nil {
				t.Fatal(err)
			}
			mutant, err := problem.ReadMutant(problemsDirForTest)
			if err != nil {
				t.Fatal(err)
			}
			hidden, err := problem.ReadTestFile(problemsDirForTest)
			if err != nil {
				t.Fatal(err)
			}
			answer, err := problem.ReadAnswer(problemsDirForTest)
			if err != nil {
				t.Fatal(err)
			}
			code := string(problems.StripBuildIgnoreTag(answer))

			stage1 := runHostGoTest(t, goMod, goSum, writesTestStageFiles(support, code, string(hidden), string(target)))
			if !stage1.passed {
				t.Fatalf("Stage 1: 模範解答が正しい実装で合格しない\n%s", stage1.output)
			}

			stage2 := runHostGoTest(t, goMod, goSum, writesTestStageFiles(support, code, string(hidden), string(mutant)))
			if stage2.passed {
				t.Fatalf("Stage 2: 模範解答のテストがバグ入り実装を検出できない\n%s", stage2.output)
			}
			if stage2.buildFailed {
				t.Fatalf("Stage 2: テスト失敗ではなくビルドエラーで落ちている（mutant.go を確認）\n%s", stage2.output)
			}
		})
	}
}

// runHostGoTest は採点と同じファイル構成のワークスペースを一時ディレクトリに作り、
// ホストの go で `go test -vet=off -json ./...` を実行する。
func runHostGoTest(t *testing.T, goMod, goSum []byte, files map[string]string) runOutcome {
	t.Helper()
	dir := t.TempDir()
	if err := writeWorkspace(dir, goMod, goSum, files); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "test", "-vet=off", "-json", "./...")
	cmd.Dir = dir
	cmd.Env = append(cmd.Environ(), "GOFLAGS=", "GOPROXY=off")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()

	return runOutcome{
		passed:      runErr == nil,
		output:      formatGoTestOutput(stdout.Bytes(), stderr.Bytes()),
		buildFailed: runErr != nil && isBuildFailure(stdout.Bytes(), stderr.Bytes()),
	}
}
