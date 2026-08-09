package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/syakesoba/codeforge/internal/auth"
	"github.com/syakesoba/codeforge/internal/problems"
	"github.com/syakesoba/codeforge/internal/store"
)

type draftRequest struct {
	Code string `json:"code"`
}

type draftResponse struct {
	Code string `json:"code"`
}

// getDraftHandler はログイン中ユーザーの自動保存済みコードを返す。
// 未保存なら 204、未ログインなら 401 を返す。
func getDraftHandler(authSvc *auth.Service, st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if _, ok := problems.Get(id); !ok {
			http.Error(w, "problem not found", http.StatusNotFound)
			return
		}

		user, err := authSvc.UserByToken(sessionToken(r))
		if err != nil {
			writeError(w, http.StatusUnauthorized, "not logged in")
			return
		}

		code, ok, err := st.GetDraft(user.ID, id)
		if err != nil {
			log.Printf("failed to get draft: %v", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		if !ok {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		writeJSON(w, http.StatusOK, draftResponse{Code: code})
	}
}

// saveDraftHandler はエディタの自動保存を受け付ける。未ログインなら 401 を返す。
func saveDraftHandler(authSvc *auth.Service, st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if _, ok := problems.Get(id); !ok {
			http.Error(w, "problem not found", http.StatusNotFound)
			return
		}

		user, err := authSvc.UserByToken(sessionToken(r))
		if err != nil {
			writeError(w, http.StatusUnauthorized, "not logged in")
			return
		}

		var req draftRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if err := st.SaveDraft(user.ID, id, req.Code); err != nil {
			log.Printf("failed to save draft: %v", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
