package main

import (
	"net"
	"net/http"

	"github.com/syakesoba/codeforge/internal/ratelimit"
)

// clientIP はレート制限のキーに使うクライアントIPを取り出す。
//
// リバースプロキシ経由での運用は未対応（X-Forwarded-For 等は見ない）。
// そのような構成にする場合は、信頼できるプロキシからのヘッダーのみを
// 参照するよう別途対応が必要。
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// rateLimited はハンドラをクライアントIP単位のレート制限でラップする。
// 制限に達した場合は 429 を返す。
func rateLimited(limiter *ratelimit.Limiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow(clientIP(r)) {
			writeError(w, http.StatusTooManyRequests, "試行回数が多すぎます。しばらく時間をおいて再度お試しください。")
			return
		}
		next(w, r)
	}
}
