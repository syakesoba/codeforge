//go:build ignore

package main

import "testing"

// gradeCase は gradeScore の1テストケースを表します。
type gradeCase struct {
	name  string
	score int
	want  string
}

// runGradeCases は cases の各要素について gradeScore を呼び出し、
// t.Run でサブテストとして実行し、結果が want と一致するか検証してください。
func runGradeCases(t *testing.T, cases []gradeCase) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := gradeScore(c.score)
			if got != c.want {
				t.Errorf("gradeScore(%d) = %q, want %q", c.score, got, c.want)
			}
		})
	}
}
