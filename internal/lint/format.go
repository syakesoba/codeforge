package lint

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/imports"
)

// undefinedPattern は `go build` の "undefined: X" 診断から識別子 X を取り出す。
var undefinedPattern = regexp.MustCompile(`^undefined: (\w+)$`)

// FixImports は goimports 相当の処理（不足importの追加・未使用importの削除・
// gofmt整形）を行います。ソースの静的解析のみでコードを実行しないため、
// Checkと同様にバックエンドプロセス内で直接実行して問題ありません。
//
// goimportsの依存解決は「対象ファイルの実際のディレクトリから一番近い
// go.mod」を辿ってモジュールの依存グラフを見に行くため、Checkと同じく
// レッスンのgo.mod/go.sumを実ファイルとして置いた一時ワークスペースの中で
// 実行する。
//
// gin等はgoimportsのヒューリスティックで解決できるが、gorm.io/gorm のように
// 見つけられないケースが確認できたため、goimports実行後に残っている
// "undefined: X" をgo.modのrequire一覧（末尾のパス要素がXと一致するもの）から
// 補完するフォールバックを追加している。
func FixImports(ctx context.Context, goMod, goSum []byte, code string) (string, error) {
	workdir, err := setupWorkspace(goMod, goSum, code)
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(workdir)

	filename := filepath.Join(workdir, "solution.go")

	formatted, err := imports.Process(filename, []byte(code), nil)
	if err != nil {
		return "", fmt.Errorf("failed to fix imports: %w", err)
	}

	formatted, err = fillKnownImports(ctx, workdir, goMod, formatted)
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}

// fillKnownImports はビルドし直して "undefined: X" が残っていないか確認し、
// 残っていればgo.modのrequire一覧から解決を試みる。
func fillKnownImports(ctx context.Context, workdir string, goModBytes, code []byte) ([]byte, error) {
	diagnostics, err := buildDiagnostics(ctx, workdir, code)
	if err != nil {
		return nil, err
	}

	modFile, err := modfile.Parse("go.mod", goModBytes, nil)
	if err != nil {
		// go.modが壊れている場合はここでは何もしない（元のコードをそのまま返す）
		return code, nil
	}

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "solution.go", code, parser.ParseComments)
	if err != nil {
		// 構文エラーがある場合はimports.Processの結果をそのまま返す
		return code, nil
	}

	added := false
	for _, d := range diagnostics {
		m := undefinedPattern.FindStringSubmatch(d.Message)
		if m == nil {
			continue
		}
		ident := m[1]

		importPath := findRequiredImportByIdent(modFile, ident)
		if importPath == "" {
			continue
		}

		if astutil.AddImport(fset, astFile, importPath) {
			added = true
		}
	}

	if !added {
		return code, nil
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, astFile); err != nil {
		return code, nil
	}

	// 追加したimportのグループ分け・並び順をgoimportsで整える
	final, err := imports.Process("solution.go", buf.Bytes(), nil)
	if err != nil {
		return buf.Bytes(), nil
	}
	return final, nil
}

// findRequiredImportByIdent は go.mod の require 一覧から、
// インポートパスの末尾の要素が ident と一致するものを探す。
// 例: "gorm.io/gorm" の末尾は "gorm"、"github.com/glebarez/sqlite" の末尾は "sqlite"。
func findRequiredImportByIdent(modFile *modfile.File, ident string) string {
	for _, req := range modFile.Require {
		if path.Base(req.Mod.Path) == ident {
			return req.Mod.Path
		}
	}
	return ""
}
