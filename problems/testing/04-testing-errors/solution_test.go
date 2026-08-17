package main

import (
	"strconv"
	"testing"
)

// TestParsePositiveInt はユーザーが実装した runParsePositiveIntCases を、
// 実際のケースで呼び出す。0以下の値を弾くErrNotPositiveのチェックを忘れた
// 実装（mutant.go）では、負の数・0のケースで不一致になり、ユーザーのテストが
// 正しくエラーを比較・失敗処理をしていれば検出できる。
func TestParsePositiveInt(t *testing.T) {
	cases := []parseCase{
		{"positive number succeeds", "42", 42, nil},
		{"one succeeds", "1", 1, nil},
		{"zero is not positive", "0", 0, ErrNotPositive},
		{"negative is not positive", "-5", 0, ErrNotPositive},
		{"non-numeric is a parse error", "abc", 0, strconv.ErrSyntax},
	}
	runParsePositiveIntCases(t, cases)
}
