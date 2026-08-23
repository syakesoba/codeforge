package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/syakesoba/codeforge/internal/auth"
	"github.com/syakesoba/codeforge/internal/store"
)

type credentialsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type progressResponse struct {
	CompletedProblemIDs []string          `json:"completedProblemIds"`
	PassedAt            map[string]string `json:"passedAt"`
}

// setSessionCookie はセッショントークンをHttpOnly Cookieとして設定します。
//
// フロントエンド(:3000)とバックエンド(:8080)はポートが異なるため、
// ブラウザからはクロスサイト扱いになります。そのため SameSite=None が必要ですが、
// SameSite=None は Secure 必須なので、ローカル開発（http）ではCookieが拒否されます。
// これを避けるため、開発時は SameSite=Lax を使います。
// 本番（HTTPS・同一ドメイン配信 or 適切なCORS設定）では Secure=true にしてください。
func setSessionCookie(w http.ResponseWriter, session *auth.Session, secure bool) {
	cookie := &http.Cookie{
		Name:     auth.SessionCookieName,
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
	if secure {
		cookie.SameSite = http.SameSiteNoneMode
	}
	http.SetCookie(w, cookie)
}

func clearSessionCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     auth.SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
	})
}

func sessionToken(r *http.Request) string {
	cookie, err := r.Cookie(auth.SessionCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode response", "err", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func signUpHandler(svc *auth.Service, secureCookie bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req credentialsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		session, err := svc.SignUp(req.Email, req.Password)
		switch {
		case errors.Is(err, auth.ErrInvalidEmail),
			errors.Is(err, auth.ErrPasswordTooShort):
			writeError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, auth.ErrEmailTaken):
			writeError(w, http.StatusConflict, "このメールアドレスは既に登録されています")
			return
		case err != nil:
			slog.Error("signup error", "err", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		setSessionCookie(w, session, secureCookie)
		writeJSON(w, http.StatusCreated, userResponse{ID: session.User.ID, Email: session.User.Email})
	}
}

func logInHandler(svc *auth.Service, secureCookie bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req credentialsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		session, err := svc.LogIn(req.Email, req.Password)
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			writeError(w, http.StatusUnauthorized, "メールアドレスまたはパスワードが正しくありません")
			return
		case err != nil:
			slog.Error("login error", "err", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		setSessionCookie(w, session, secureCookie)
		writeJSON(w, http.StatusOK, userResponse{ID: session.User.ID, Email: session.User.Email})
	}
}

func logOutHandler(svc *auth.Service, secureCookie bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := svc.LogOut(sessionToken(r)); err != nil {
			slog.Error("logout error", "err", err)
		}
		clearSessionCookie(w, secureCookie)
		w.WriteHeader(http.StatusNoContent)
	}
}

func meHandler(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := svc.UserByToken(sessionToken(r))
		if err != nil {
			writeError(w, http.StatusUnauthorized, "not logged in")
			return
		}
		writeJSON(w, http.StatusOK, userResponse{ID: user.ID, Email: user.Email})
	}
}

func progressHandler(svc *auth.Service, st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := svc.UserByToken(sessionToken(r))
		if err != nil {
			writeError(w, http.StatusUnauthorized, "not logged in")
			return
		}

		entries, err := st.CompletedProgress(user.ID)
		if err != nil {
			slog.Error("failed to list progress", "err", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		ids := make([]string, len(entries))
		passedAt := make(map[string]string, len(entries))
		for i, e := range entries {
			ids[i] = e.ProblemID
			passedAt[e.ProblemID] = e.PassedAt.Format(time.RFC3339)
		}

		writeJSON(w, http.StatusOK, progressResponse{CompletedProblemIDs: ids, PassedAt: passedAt})
	}
}
