package main

import "testing"

// TestGradeScore はユーザーが実装した runGradeCases を、実際のケースで呼び出す。
// 採点は2段階で行われる:
//  1. 正しい実装（target.go）に対してこのテストが合格するか
//  2. わざとバグを仕込んだ実装（mutant.go）に対してこのテストが不合格になるか
//
// 両方を満たして初めて合格になるため、比較や失敗処理を書かない空の実装では
// 合格できない（内部の詳細は internal/judge の2段階採点を参照）。
func TestGradeScore(t *testing.T) {
	cases := []gradeCase{
		{"perfect score is A", 100, "A"},
		{"boundary of A", 90, "A"},
		{"just below A is B", 89, "B"},
		{"boundary of B", 70, "B"},
		{"just below B is C", 69, "C"},
		{"boundary of C", 50, "C"},
		{"just below C is F", 49, "F"},
		{"zero is F", 0, "F"},
	}
	runGradeCases(t, cases)
}
