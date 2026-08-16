package main

import "testing"

// TestDaysInMonth はユーザーが実装した runDaysInMonthCases を、実際のケースで呼び出す。
// うるう年の2月（29日になるケース）を含めているため、うるう年判定を忘れた
// 実装（mutant.go）ではこのケースで不一致になり、ユーザーのテストが正しく
// 比較・失敗処理をしていれば検出できる。
func TestDaysInMonth(t *testing.T) {
	cases := []daysCase{
		{1, false, 31},
		{2, false, 28},
		{2, true, 29},
		{4, false, 30},
		{4, true, 30},
		{12, false, 31},
		{12, true, 31},
	}
	runDaysInMonthCases(t, cases)
}
