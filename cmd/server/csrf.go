package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"log"
	"net/http"
	"time"
)

// csrfCookieName / csrfHeaderName は Double Submit Cookie 方式のCSRF対策に使う。
//
// セッションCookieは SameSite=None（本番）でクロスオリジンでも送信されるため、
// SameSite単独ではCSRFを防げない。そこでJavaScriptから読める非HttpOnlyの
// トークンをCookieとして発行し、状態変更を伴うリクエストにはこれを
// 独自ヘッダー（X-CSRF-Token）にも載せて送らせる。攻撃者のサイトは
// 同一生成元ポリシーによりこのCookieの値を読めないため、ヘッダーを
// 正しく複製できず、フォームの自動送信等によるCSRFを防止できる。
const (
	csrfCookieName = "csrf_token"
	csrfHeaderName = "X-CSRF-Token"
)

func generateCSRFToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// ensureCSRFCookie は既存のCSRFトークンCookieが無ければ新規発行する。
// 全リクエストに対して呼ぶことで、状態変更リクエストを送る前に
// フロントエンドが必ずトークンを取得できている状態にする。
func ensureCSRFCookie(w http.ResponseWriter, r *http.Request, secure bool) {
	if _, err := r.Cookie(csrfCookieName); err == nil {
		return
	}

	token, err := generateCSRFToken()
	if err != nil {
		log.Printf("failed to generate csrf token: %v", err)
		return
	}

	cookie := &http.Cookie{
		Name:    csrfCookieName,
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(24 * time.Hour),
		// JavaScriptから読み取ってヘッダーに載せる必要があるため HttpOnly にはしない。
		HttpOnly: false,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
	if secure {
		cookie.SameSite = http.SameSiteNoneMode
	}
	http.SetCookie(w, cookie)
}

// verifyCSRF はCookieとヘッダーのトークンが一致するかを検証する。
func verifyCSRF(r *http.Request) bool {
	cookie, err := r.Cookie(csrfCookieName)
	if err != nil || cookie.Value == "" {
		return false
	}
	header := r.Header.Get(csrfHeaderName)
	if header == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(header)) == 1
}

// withCSRFCookie は全リクエストにCSRFトークンCookieを発行するミドルウェア。
// 検証自体は requireCSRF を個別のハンドラに適用して行う
// （/check, /format のような状態を変更しないエンドポイントまで検証すると、
// クッキーを送らないリクエストで不要に失敗するため）。
func withCSRFCookie(secure bool, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ensureCSRFCookie(w, r, secure)
		h.ServeHTTP(w, r)
	})
}

// requireCSRF は状態変更を伴うハンドラに適用し、トークンの一致を必須にする。
func requireCSRF(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !verifyCSRF(r) {
			writeError(w, http.StatusForbidden, "CSRFトークンが無効です。ページを再読み込みしてください。")
			return
		}
		next(w, r)
	}
}
