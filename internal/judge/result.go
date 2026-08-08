package judge

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"
)

// testEvent は `go test -json` が1行ごとに出力するイベントです。
// 詳細: https://pkg.go.dev/cmd/test2json
type testEvent struct {
	Action string
	Test   string
	Output string
}

const maxOutputBytes = 32 * 1024

// formatGoTestOutput は `go test -json` の標準出力を解析し、
// 人間が読める実行ログを組み立てる。合否判定には使わない
// （judge.Run で go test プロセスの終了コードを正として扱う）。
//
// ビルドエラーなど、テストバイナリが生成される前に失敗した場合はJSONイベントが
// 一切出力されないため、その場合は stderr の内容をそのまま出力として扱う。
func formatGoTestOutput(stdout, stderr []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(stdout))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var logLines strings.Builder
	sawEvent := false

	for scanner.Scan() {
		var ev testEvent
		line := scanner.Bytes()
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		if err := json.Unmarshal(line, &ev); err != nil {
			// JSONとして解釈できない行はそのままログに含める
			logLines.Write(line)
			logLines.WriteByte('\n')
			continue
		}
		sawEvent = true
		if ev.Action == "output" {
			logLines.WriteString(ev.Output)
		}
	}

	if !sawEvent {
		// テストが1件も実行されていない = コンパイルエラーなど致命的な失敗
		combined := strings.TrimSpace(string(stderr))
		if combined == "" {
			combined = strings.TrimSpace(string(stdout))
		}
		return truncate(combined)
	}

	return truncate(strings.TrimSpace(logLines.String()))
}

func truncate(s string) string {
	if len(s) <= maxOutputBytes {
		return s
	}
	return s[:maxOutputBytes] + "\n...(output truncated)"
}
