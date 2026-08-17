package main

import (
	"io"
	"log/slog"
)

// NewRequestLogger は、標準出力（コンテナ環境では docker logs で収集される）などの
// io.Writer に対して、JSON形式でログを出力する *slog.Logger を作って返します。
func NewRequestLogger(w io.Writer) *slog.Logger {
	// TODO: slog.NewJSONHandler(w, nil) でハンドラーを作り、
	// slog.New(handler) でLoggerを作って返す
	return nil
}

// LogRequest は、logger を使ってHTTPリクエストの処理結果を1件、
// Infoレベルのログとして出力します。
// メッセージは "request handled"、属性として method, path, status を付与してください。
func LogRequest(logger *slog.Logger, method, path string, status int) {
	// TODO: logger.Info("request handled", "method", method, "path", path, "status", status) を呼ぶ
}

func main() {}
