//go:build ignore

package main

import (
	"testing"
	"time"
)

// fakeClock はテスト用の Clock 実装です。
type fakeClock struct {
	now time.Time
}

func (f fakeClock) Now() time.Time {
	return f.now
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
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			clock := fakeClock{now: c.now}
			got := isExpired(clock, c.expiresAt)
			if got != c.want {
				t.Errorf("isExpired(now=%v, expiresAt=%v) = %v, want %v", c.now, c.expiresAt, got, c.want)
			}
		})
	}
}
