package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandlerHealthy(t *testing.T) {
	handler := HealthHandler(func() error { return nil })

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != "ok" {
		t.Errorf("body = %q, want %q", rec.Body.String(), "ok")
	}
}

func TestHealthHandlerUnhealthy(t *testing.T) {
	handler := HealthHandler(func() error { return errors.New("database unreachable") })

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
	if rec.Body.String() != "database unreachable" {
		t.Errorf("body = %q, want %q", rec.Body.String(), "database unreachable")
	}
}

func TestHealthHandlerCallsCheckerEachTime(t *testing.T) {
	calls := 0
	handler := HealthHandler(func() error {
		calls++
		return nil
	})

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		handler(rec, req)
	}

	if calls != 3 {
		t.Errorf("checkが呼ばれた回数 = %d, want 3", calls)
	}
}
