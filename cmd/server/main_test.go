package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	t.Run("未設定なら既定値", func(t *testing.T) {
		cfg := loadConfig(func(string) string { return "" })
		want := config{
			Port:            "8080",
			DBPath:          "data/codeforge.db",
			FrontendOrigin:  "http://localhost:3000",
			SecureCookie:    false,
			ShutdownTimeout: 10 * time.Second,
			JudgeImage:      "codeforge-judge:latest",
			JudgeWorkDir:    "",
		}
		if cfg != want {
			t.Errorf("loadConfig() = %+v, want %+v", cfg, want)
		}
	})

	t.Run("環境変数で上書きできる", func(t *testing.T) {
		env := map[string]string{
			"PORT":                     "9090",
			"DB_PATH":                  "/data/app.db",
			"FRONTEND_ORIGIN":          "https://example.com",
			"SECURE_COOKIE":            "true",
			"SHUTDOWN_TIMEOUT_SECONDS": "30",
			"JUDGE_IMAGE":              "registry.example.com/judge:v2",
			"JUDGE_WORK_DIR":           "/tmp/codeforge-judge",
		}
		cfg := loadConfig(func(k string) string { return env[k] })
		want := config{
			Port:            "9090",
			DBPath:          "/data/app.db",
			FrontendOrigin:  "https://example.com",
			SecureCookie:    true,
			ShutdownTimeout: 30 * time.Second,
			JudgeImage:      "registry.example.com/judge:v2",
			JudgeWorkDir:    "/tmp/codeforge-judge",
		}
		if cfg != want {
			t.Errorf("loadConfig() = %+v, want %+v", cfg, want)
		}
	})
}

func TestReadyHandler(t *testing.T) {
	ok := func(context.Context) error { return nil }
	ng := func(context.Context) error { return errors.New("docker command not found") }

	tests := []struct {
		name       string
		checks     []readinessCheck
		wantStatus int
		wantBody   map[string]string
	}{
		{
			name:       "全て成功なら200",
			checks:     []readinessCheck{{"db", ok}, {"judge", ok}, {"lint", ok}},
			wantStatus: http.StatusOK,
			wantBody:   map[string]string{"db": "ok", "judge": "ok", "lint": "ok"},
		},
		{
			name:       "1つでも失敗なら503とエラー内容",
			checks:     []readinessCheck{{"db", ok}, {"judge", ng}, {"lint", ok}},
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   map[string]string{"db": "ok", "judge": "docker command not found", "lint": "ok"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			readyHandler(tt.checks...)(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			var got map[string]string
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("invalid JSON: %v: %s", err, rec.Body.String())
			}
			if len(got) != len(tt.wantBody) {
				t.Fatalf("body = %v, want %v", got, tt.wantBody)
			}
			for k, v := range tt.wantBody {
				if got[k] != v {
					t.Errorf("body[%q] = %q, want %q", k, got[k], v)
				}
			}
		})
	}
}
