package main

import "testing"

// ErrNotPositive・parsePositiveInt は別ファイル（target.go）で実装済みです。

// parseCase は parsePositiveInt の1テストケースを表します。
// wantErrIs が nil なら成功（wantの値が返る）を期待し、
// nilでなければ errors.Is(err, wantErrIs) が true になることを期待します。
type parseCase struct {
	name      string
	input     string
	want      int
	wantErrIs error
}

// runParsePositiveIntCases は cases の各要素について parsePositiveInt を呼び出し、
// t.Run でサブテストとして実行し、値またはエラーが期待通りか検証してください。
func runParsePositiveIntCases(t *testing.T, cases []parseCase) {
	// TODO:
	// 1. cases をfor-rangeでループし、t.Run(c.name, func(t *testing.T) { ... }) でサブテストを作る
	// 2. got, err := parsePositiveInt(c.input) を呼び出す
	// 3. c.wantErrIs が nil でない場合:
	//    errors.Is(err, c.wantErrIs) が false なら t.Errorf で失敗させる
	// 4. c.wantErrIs が nil の場合:
	//    err が nil でなければ t.Errorf で失敗させる
	//    got が c.want と一致しなければ t.Errorf で失敗させる
	//
	// errors.Is を使うには "errors" パッケージのimportが必要です。
	// エディタの「インポートを自動修正」機能で追加できます。
}
