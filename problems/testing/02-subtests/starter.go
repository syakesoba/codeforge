package main

import (
	"fmt"
	"testing"
)

// daysInMonth は別ファイル（target.go）で実装済みです。
// 月（1〜12）とうるう年かどうかを受け取り、その月の日数を返します。

// daysCase は daysInMonth の1テストケースを表します。
// Lesson 1のgradeCaseと違い、nameフィールドがありません。
type daysCase struct {
	month      int
	isLeapYear bool
	want       int
}

// runDaysInMonthCases は cases の各要素について daysInMonth を呼び出し、
// t.Run でサブテストとして実行し、結果が want と一致するか検証してください。
//
// nameフィールドが無いので、fmt.Sprintf で入力値からサブテストの名前を
// 組み立ててください（例: "month=2/leap=true"）。
func runDaysInMonthCases(t *testing.T, cases []daysCase) {
	// TODO:
	// 1. cases をfor-rangeでループする
	// 2. fmt.Sprintf("month=%d/leap=%v", c.month, c.isLeapYear) でサブテスト名を作る
	// 3. t.Run(name, func(t *testing.T) { ... }) でサブテストを作る
	// 4. daysInMonth(c.month, c.isLeapYear) の結果が c.want と一致するか確認し、
	//    一致しなければ t.Errorf で失敗させる
}
