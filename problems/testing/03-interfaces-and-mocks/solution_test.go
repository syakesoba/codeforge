package main

import (
	"testing"
	"time"
)

// TestIsExpired はユーザーが実装した runIsExpiredCases を、実際のケースで呼び出す。
// 「期限前」「期限後」の両方を含めているため、AfterとBeforeを取り違えた実装
// （mutant.go）ではいずれかのケースで不一致になり、ユーザーのテストが正しく
// 比較・失敗処理をしていれば検出できる。
func TestIsExpired(t *testing.T) {
	base := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	cases := []expiredCase{
		{"before expiry is not expired", base, base.Add(time.Hour), false},
		{"exactly at expiry is not expired", base, base, false},
		{"after expiry is expired", base, base.Add(-time.Hour), true},
	}
	runIsExpiredCases(t, cases)
}
