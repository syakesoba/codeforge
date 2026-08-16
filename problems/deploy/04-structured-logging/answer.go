//go:build ignore

package main

import (
	"io"
	"log/slog"
)

func NewRequestLogger(w io.Writer) *slog.Logger {
	handler := slog.NewJSONHandler(w, nil)
	return slog.New(handler)
}

func LogRequest(logger *slog.Logger, method, path string, status int) {
	logger.Info("request handled", "method", method, "path", path, "status", status)
}

func main() {}
