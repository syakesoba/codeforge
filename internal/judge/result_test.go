package judge

import "testing"

func TestIsBuildFailure(t *testing.T) {
	tests := []struct {
		name   string
		stdout string
		stderr string
		want   bool
	}{
		{
			name: "Go 1.24以降のbuild-failイベント",
			stdout: `{"ImportPath":"submission [submission.test]","Action":"build-output","Output":"undefined: gradeScore\n"}
{"ImportPath":"submission [submission.test]","Action":"build-fail"}
{"Action":"output","Package":"submission","Output":"FAIL\tsubmission [build failed]\n"}
{"Action":"fail","Package":"submission","Elapsed":0,"FailedBuild":"submission [submission.test]"}`,
			want: true,
		},
		{
			name:   "古い形式（JSONイベントが無くstdoutに定型文だけ出る）",
			stdout: "FAIL\tsubmission [build failed]\n",
			stderr: "./solution_test.go:18:11: undefined: gradeScore\n",
			want:   true,
		},
		{
			name: "テストの失敗はビルド失敗ではない",
			stdout: `{"Action":"run","Package":"submission","Test":"TestGradeScore"}
{"Action":"output","Package":"submission","Test":"TestGradeScore","Output":"gradeScore(90) = \"B\", want \"A\"\n"}
{"Action":"fail","Package":"submission","Test":"TestGradeScore"}
{"Action":"fail","Package":"submission"}`,
			want: false,
		},
		{
			name:   "合格",
			stdout: `{"Action":"pass","Package":"submission"}`,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBuildFailure([]byte(tt.stdout), []byte(tt.stderr)); got != tt.want {
				t.Errorf("isBuildFailure() = %v, want %v", got, tt.want)
			}
		})
	}
}
