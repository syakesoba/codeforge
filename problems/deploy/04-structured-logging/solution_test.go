package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewRequestLoggerProducesJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := NewRequestLogger(&buf)
	if logger == nil {
		t.Fatal("NewRequestLoggerがnilを返しました")
	}

	LogRequest(logger, "GET", "/healthz", 200)

	line := strings.TrimSpace(buf.String())
	if line == "" {
		t.Fatal("ログが出力されていません")
	}

	var entry map[string]any
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		t.Fatalf("出力がJSONとしてパースできません: %v\n出力: %s", err, line)
	}
}

func TestLogRequestFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewRequestLogger(&buf)

	LogRequest(logger, "POST", "/api/orders", 201)

	var entry map[string]any
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry); err != nil {
		t.Fatalf("出力がJSONとしてパースできません: %v", err)
	}

	if entry["msg"] != "request handled" {
		t.Errorf("msg = %v, want %q", entry["msg"], "request handled")
	}
	if entry["method"] != "POST" {
		t.Errorf("method = %v, want %q", entry["method"], "POST")
	}
	if entry["path"] != "/api/orders" {
		t.Errorf("path = %v, want %q", entry["path"], "/api/orders")
	}
	// JSONの数値はfloat64としてデコードされる
	status, ok := entry["status"].(float64)
	if !ok || status != 201 {
		t.Errorf("status = %v, want 201", entry["status"])
	}
}

func TestLogRequestEachCallIsOneLine(t *testing.T) {
	var buf bytes.Buffer
	logger := NewRequestLogger(&buf)

	LogRequest(logger, "GET", "/a", 200)
	LogRequest(logger, "GET", "/b", 404)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("出力行数 = %d, want 2（1回のログ出力が1行のJSONになっているか確認してください）", len(lines))
	}

	var second map[string]any
	if err := json.Unmarshal([]byte(lines[1]), &second); err != nil {
		t.Fatalf("2行目がJSONとしてパースできません: %v", err)
	}
	if second["path"] != "/b" {
		t.Errorf("2行目のpath = %v, want %q", second["path"], "/b")
	}
}
