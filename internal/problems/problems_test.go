package problems

import (
	"bytes"
	"testing"
)

func TestStripBuildIgnoreTag(t *testing.T) {
	tests := []struct {
		name string
		src  string
		want string
	}{
		{"タグと空行を取り除く", "//go:build ignore\n\npackage main\n", "package main\n"},
		{"CRLF改行でも取り除く", "//go:build ignore\r\n\r\npackage main\r\n", "package main\r\n"},
		{"直後が空行でなければタグ行だけ取り除く", "//go:build ignore\npackage main\n", "package main\n"},
		{"タグが無ければそのまま", "package main\n", "package main\n"},
		{"先頭行以外のタグは取り除かない", "package main\n//go:build ignore\n", "package main\n//go:build ignore\n"},
		{"別のビルド制約はそのまま", "//go:build linux\n\npackage main\n", "//go:build linux\n\npackage main\n"},
		{"タグだけのファイルは空になる", "//go:build ignore", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(StripBuildIgnoreTag([]byte(tt.src)))
			if got != tt.want {
				t.Errorf("StripBuildIgnoreTag(%q) = %q, want %q", tt.src, got, tt.want)
			}
		})
	}
}

// TestWritesTestFilesAreBuildable は、WritesTest レッスンの target.go / mutant.go を
// 読み込んだ結果に `//go:build ignore` が残っていないことを確認する。
// 残っていると採点ワークスペースでそのファイルがビルドから除外され、
// 2段階採点が正しく働かない。
func TestWritesTestFilesAreBuildable(t *testing.T) {
	const baseDir = "../../problems"
	count := 0
	for id, p := range allowlist {
		if !p.WritesTest {
			continue
		}
		count++
		for name, read := range map[string]func(string) ([]byte, error){
			"target.go": p.ReadTarget,
			"mutant.go": p.ReadMutant,
		} {
			b, err := read(baseDir)
			if err != nil {
				t.Errorf("%s: %s を読み込めない: %v", id, name, err)
				continue
			}
			if bytes.Contains(b, []byte(buildIgnoreTag)) {
				t.Errorf("%s: 読み込んだ %s にビルド制約 %q が残っている", id, name, buildIgnoreTag)
			}
		}
	}
	if count == 0 {
		t.Fatal("WritesTest のレッスンが1件も無い")
	}
}
