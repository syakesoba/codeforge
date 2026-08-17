//go:build ignore

package main

import "time"

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// isExpired はわざとバグを仕込んだ実装です
// （AfterではなくBeforeを使っており、判定が逆になってしまう）。
func isExpired(clock Clock, expiresAt time.Time) bool {
	return clock.Now().Before(expiresAt)
}
