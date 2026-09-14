package judge

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/syakesoba/codeforge/internal/problems"
)

// dockerTestsEnv を "1" にしたときだけ、実際の codeforge-judge イメージを使う
// 採点テストを実行する（Docker とイメージのビルドが必要なため、通常の
// `go test ./...` ではスキップする）。
const dockerTestsEnv = "CODEFORGE_DOCKER_TESTS"

// problemsDirForTest は、テスト実行時のカレントディレクトリ（internal/judge）から見た
// problems/ ディレクトリの位置。
const problemsDirForTest = "../../problems"

// assertionlessTesting01 は testing-01 の非公開テストが要求する型・関数だけを
// 定義し、結果の比較（t.Errorf）を一切しない「検証していないテスト」。
// 2段階採点が正しく働いていれば、Stage 2（バグ入り実装）でも合格してしまうため
// insufficient_test として不合格になるはず。
const assertionlessTesting01 = `package main

import "testing"

type gradeCase struct {
	name  string
	score int
	want  string
}

func runGradeCases(t *testing.T, cases []gradeCase) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_ = gradeScore(c.score)
		})
	}
}
`

func TestRunWritesTestWithDocker(t *testing.T) {
	if os.Getenv(dockerTestsEnv) != "1" {
		t.Skipf("%s=1 のときだけ実行する（Docker と codeforge-judge イメージが必要）", dockerTestsEnv)
	}

	problem, ok := problems.Get("testing-01")
	if !ok {
		t.Fatal("testing-01 が許可リストに無い")
	}
	answer, err := problem.ReadAnswer(problemsDirForTest)
	if err != nil {
		t.Fatal(err)
	}

	runner := NewRunner()
	runner.ProblemsBaseDir = problemsDirForTest

	tests := []struct {
		name       string
		code       string
		wantPassed bool
		wantError  string
	}{
		{
			name:       "模範解答のテストは合格する",
			code:       string(problems.StripBuildIgnoreTag(answer)),
			wantPassed: true,
		},
		{
			name:       "比較をしないテストはバグを検出できないので不合格",
			code:       assertionlessTesting01,
			wantPassed: false,
			wantError:  "insufficient_test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()

			got, err := runner.Run(ctx, problem, tt.code)
			if err != nil {
				t.Fatalf("Run returned error: %v", err)
			}
			if got.Passed != tt.wantPassed || got.Error != tt.wantError {
				t.Errorf("Passed=%v Error=%q, want Passed=%v Error=%q\noutput:\n%s",
					got.Passed, got.Error, tt.wantPassed, tt.wantError, got.Output)
			}
		})
	}
}
