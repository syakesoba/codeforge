package main

import "testing"

// gradeScore は別ファイル（target.go）で実装済みです。
// score >= 90 なら "A"、70以上なら "B"、50以上なら "C"、それ未満なら "F" を返します。

// gradeCase は gradeScore の1テストケースを表します。
type gradeCase struct {
	name  string
	score int
	want  string
}

// runGradeCases は cases の各要素について gradeScore を呼び出し、
// t.Run でサブテストとして実行し、結果が want と一致するか検証してください。
func runGradeCases(t *testing.T, cases []gradeCase) {
	// TODO:
	// 1. cases をfor-rangeでループする
	// 2. 各ケースについて t.Run(c.name, func(t *testing.T) { ... }) でサブテストを作る
	// 3. サブテストの中で gradeScore(c.score) を呼び、結果が c.want と一致するか確認する
	// 4. 一致しなければ t.Errorf で失敗させる（メッセージに実際の値と期待値を含めること）
}
