package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// greetHandler は別ファイル（target.go）で実装済みです。

// greetCase は greetHandler の1テストケースを表します。
type greetCase struct {
	name       string
	query      string // 例: "?name=Alice"、無しなら ""
	wantStatus int
	wantBody   string // wantStatusが200のときだけ比較する
}

// runGreetHandlerCases は cases の各要素について、httptest でリクエストを作り
// greetHandler を呼び出し、ステータスコード・本文が期待通りか検証してください。
func runGreetHandlerCases(t *testing.T, cases []greetCase) {
	// TODO:
	// 1. cases をfor-rangeでループし、t.Run(c.name, func(t *testing.T) { ... }) でサブテストを作る
	// 2. req := httptest.NewRequest(http.MethodGet, "/greet"+c.query, nil) でリクエストを作る
	// 3. rec := httptest.NewRecorder() でレスポンスを記録するレコーダーを作る
	// 4. greetHandler(rec, req) を呼び出す（普通のhttp.HandlerFuncとして直接呼べる）
	// 5. rec.Code が c.wantStatus と一致するか確認する（不一致ならt.Errorf）
	// 6. c.wantStatus が http.StatusOK の場合のみ、rec.Body.String() が
	//    c.wantBody と一致するか確認する
}
