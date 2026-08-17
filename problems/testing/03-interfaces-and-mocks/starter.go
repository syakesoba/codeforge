package main

import (
	"testing"
	"time"
)

// Clock・realClock・isExpired は別ファイル（target.go）で実装済みです。

// fakeClock はテスト用の Clock 実装です。now フィールドに保持した
// 固定時刻を返すことで、「現在時刻」をテストの中でコントロールできます。
type fakeClock struct {
	now time.Time
}

// Now は f.now をそのまま返してください（Clockインターフェースの実装）。
func (f fakeClock) Now() time.Time {
	// TODO: f.now を返す
	return time.Time{}
}

// expiredCase は isExpired の1テストケースを表します。
type expiredCase struct {
	name      string
	now       time.Time
	expiresAt time.Time
	want      bool
}

// runIsExpiredCases は cases の各要素について、その now を持つ fakeClock を作り、
// isExpired(clock, c.expiresAt) を呼び出して結果を検証してください。
func runIsExpiredCases(t *testing.T, cases []expiredCase) {
	// TODO:
	// 1. cases をfor-rangeでループする
	// 2. t.Run(c.name, func(t *testing.T) { ... }) でサブテストを作る
	// 3. clock := fakeClock{now: c.now} を作る
	// 4. isExpired(clock, c.expiresAt) の結果が c.want と一致するか確認し、
	//    一致しなければ t.Errorf で失敗させる
}
