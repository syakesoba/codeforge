package main

import "time"

// Clock は現在時刻を返すインターフェースです（実装済み・変更不要）。
// 本番では realClock（time.Now()を返すだけ）を使いますが、テストでは
// このインターフェースを実装した「偽物（フェイク）」に差し替えることで、
// 「現在時刻」に依存する処理を再現性よくテストできます。
type Clock interface {
	Now() time.Time
}

// realClock は本番で使う Clock の実装です。
type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// isExpired は、clock.Now() が expiresAt より後（after）であれば true を返します。
func isExpired(clock Clock, expiresAt time.Time) bool {
	return clock.Now().After(expiresAt)
}
