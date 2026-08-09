// Package ratelimit はキー（IPアドレス等）ごとの固定ウィンドウ方式レート制限を提供します。
// 単一プロセスでの運用を前提としたインメモリ実装で、外部ストア（Redis等）は使いません。
package ratelimit

import (
	"sync"
	"time"
)

// Limiter はキーごとに「window時間あたりmax回まで」を許可するレート制限器です。
type Limiter struct {
	mu       sync.Mutex
	max      int
	window   time.Duration
	counters map[string]*counterEntry
}

type counterEntry struct {
	count     int
	windowEnd time.Time
}

// cleanupInterval はメモリリークを防ぐため、期限切れエントリを定期的に掃除する間隔です。
const cleanupInterval = 10 * time.Minute

// New はキーごとに window 時間あたり max 回までを許可する Limiter を作成します。
// バックグラウンドで期限切れエントリを掃除するgoroutineを起動します
// （プロセス終了まで動き続ける想定のため、Close等の停止手段は設けていません）。
func New(max int, window time.Duration) *Limiter {
	l := &Limiter{
		max:      max,
		window:   window,
		counters: make(map[string]*counterEntry),
	}
	go l.cleanupLoop()
	return l
}

// Allow はキーがまだ制限内であれば true を返し、カウントを1つ消費します。
// ウィンドウを過ぎていれば自動的にリセットされます。
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	e, ok := l.counters[key]
	if !ok || now.After(e.windowEnd) {
		e = &counterEntry{count: 0, windowEnd: now.Add(l.window)}
		l.counters[key] = e
	}

	if e.count >= l.max {
		return false
	}
	e.count++
	return true
}

func (l *Limiter) cleanupLoop() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		l.mu.Lock()
		for key, e := range l.counters {
			if now.After(e.windowEnd) {
				delete(l.counters, key)
			}
		}
		l.mu.Unlock()
	}
}
