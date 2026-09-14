package problems

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/mod/modfile"
)

// primingDir は採点イメージにモジュールキャッシュを焼き込むためのモジュール
// （judge-image/priming）の、このパッケージから見た位置。
const primingDir = "../../judge-image/priming"

// TestExternalModuleLessonsMatchPriming は、外部モジュールを使うレッスン（go.sum を持つ
// レッスン）の go.mod / go.sum が、採点イメージの priming モジュールと一致していることを
// 確認する。
//
// 採点コンテナは --network none かつ GOPROXY=off で動くため、priming で焼き込んだのと
// 違うバージョンを1つでも要求すると「module lookup disabled by GOPROXY=off」や
// 「[setup failed]」で模範解答すら採点できなくなる（gRPCコース追加で priming の依存
// バージョンが上がった際、既存12レッスンがこの状態になっていた）。
func TestExternalModuleLessonsMatchPriming(t *testing.T) {
	primingModBytes, err := os.ReadFile(filepath.Join(primingDir, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	primingSum, err := os.ReadFile(filepath.Join(primingDir, "go.sum"))
	if err != nil {
		t.Fatal(err)
	}
	primingRequires := requireVersions(t, "priming go.mod", primingModBytes)

	const baseDir = "../../problems"
	checked := 0
	for id, p := range allowlist {
		sum, err := p.ReadGoSum(baseDir)
		if errors.Is(err, os.ErrNotExist) {
			continue // 標準ライブラリだけのレッスン
		}
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		checked++

		if !bytes.Equal(sum, primingSum) {
			t.Errorf("%s: go.sum が judge-image/priming/go.sum と一致しない", id)
		}

		modBytes, err := p.ReadGoMod(baseDir)
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		for path, want := range requireVersions(t, id, modBytes) {
			got, ok := primingRequires[path]
			if !ok {
				t.Errorf("%s: %s は priming の go.mod に無い（採点イメージにキャッシュされない）", id, path)
				continue
			}
			if got != want {
				t.Errorf("%s: %s のバージョンが priming と違う（レッスン %s / priming %s）", id, path, want, got)
			}
		}
	}
	if checked == 0 {
		t.Fatal("外部モジュールを使うレッスンが1件も無い")
	}
}

// requireVersions は go.mod の require（indirect含む）をモジュールパス→バージョンの
// マップにして返す。
func requireVersions(t *testing.T, name string, goMod []byte) map[string]string {
	t.Helper()
	f, err := modfile.Parse(name, goMod, nil)
	if err != nil {
		t.Fatalf("%s: go.mod を解析できない: %v", name, err)
	}
	versions := make(map[string]string, len(f.Require))
	for _, r := range f.Require {
		versions[r.Mod.Path] = r.Mod.Version
	}
	return versions
}
