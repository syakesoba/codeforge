package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewAppHealthz(t *testing.T) {
	app := NewApp(Config{}, func() error { return nil }, nil)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestNewAppCustomRoutes(t *testing.T) {
	routes := map[string]http.HandlerFunc{
		"/hello": func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hello"))
		},
		"/bye": func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("bye"))
		},
	}
	app := NewApp(Config{}, func() error { return nil }, routes)

	for path, want := range map[string]string{"/hello": "hello", "/bye": "bye"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		app.ServeHTTP(rec, req)
		if rec.Body.String() != want {
			t.Errorf("%s のレスポンス = %q, want %q", path, rec.Body.String(), want)
		}
	}
}

func TestNewAppDebugDisabledByDefault(t *testing.T) {
	app := NewApp(Config{Debug: false}, func() error { return nil }, nil)

	req := httptest.NewRequest(http.MethodGet, "/debug", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code == http.StatusOK {
		t.Error("Debug: falseなのに/debugが200を返しています（未登録のパスとして404になるべきです）")
	}
}

func TestNewAppDebugEnabled(t *testing.T) {
	app := NewApp(Config{Debug: true}, func() error { return nil }, nil)

	req := httptest.NewRequest(http.MethodGet, "/debug", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "debug mode enabled" {
		t.Errorf("body = %q, want %q", rec.Body.String(), "debug mode enabled")
	}
}

func TestNewAppHealthzReflectsChecker(t *testing.T) {
	app := NewApp(Config{}, func() error { return http.ErrBodyNotAllowed }, nil)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want %d（HealthCheckerがエラーを返す場合）", rec.Code, http.StatusServiceUnavailable)
	}
}
