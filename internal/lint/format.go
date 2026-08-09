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
	"strings"

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
//
// hints には、呼び出し側（エディタの/checkの結果）で既に判明している未解決の
// 識別子名を渡せる。指定された場合は go build による再検証（fillKnownImports）を
// 省略できるため応答が大幅に速くなる。
//
// さらに、hints解決分のimportは imports.Process を呼ぶ前にASTへ直接注入する
// （injectHintedImports）。gorm.io/gorm のようにgoimportsが解決できないパッケージを
// 先に解決しておくことで、imports.Process が「未知の識別子を探して失敗する」という
// 一番コストの高い処理を避けられる。残りのstdlib解決や整形は最後の
// imports.Process 1回にまとめて任せる（hints無しの場合は素直に1回呼ぶだけ）。
func FixImports(ctx context.Context, goMod, goSum []byte, code string, hints []string) (string, error) {
	workdir, err := setupWorkspace(goMod, goSum, code)
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(workdir)

	filename := filepath.Join(workdir, "solution.go")
	src := []byte(code)

	if len(hints) > 0 {
		src = injectHintedImports(src, goMod, hints)
	}

	formatted, err := imports.Process(filename, src, nil)
	if err != nil {
		return "", fmt.Errorf("failed to fix imports: %w", err)
	}

	if len(hints) == 0 {
		formatted, err = fillKnownImports(ctx, workdir, goMod, formatted)
		if err != nil {
			return "", err
		}
	}

	return string(formatted), nil
}

// injectHintedImports は、hints に含まれる識別子について go.mod の require 一覧
// から直接解決し、ASTレベルで先に import 文を追加する。go.modのパースやASTの
// 解析に失敗した場合、あるいは追加すべきものが無い場合は元のコードをそのまま返す
// （エラーはimports.Process側で改めて検知させる）。
func injectHintedImports(code []byte, goModBytes []byte, hints []string) []byte {
	modFile, err := modfile.Parse("go.mod", goModBytes, nil)
	if err != nil {
		return code
	}

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "solution.go", code, parser.ParseComments)
	if err != nil {
		// 構文エラーがある場合はここでは何もしない
		return code
	}

	// 既にimport済みの識別子（エイリアスがあればエイリアス、無ければパス末尾）は
	// スキップする。
	resolved := make(map[string]bool)
	for _, imp := range astFile.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		name := path.Base(importPath)
		if imp.Name != nil {
			name = imp.Name.Name
		}
		resolved[name] = true
	}

	added := false
	for _, hint := range hints {
		if resolved[hint] {
			continue
		}
		importPath := findRequiredImportByIdent(modFile, hint)
		if importPath == "" {
			continue
		}
		if astutil.AddImport(fset, astFile, importPath) {
			added = true
			resolved[hint] = true
		}
	}

	if !added {
		return code
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, astFile); err != nil {
		return code
	}
	return buf.Bytes()
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
