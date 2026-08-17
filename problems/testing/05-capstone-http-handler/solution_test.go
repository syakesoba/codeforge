package main

import (
	"net/http"
	"testing"
)

// TestGreetHandler はユーザーが実装した runGreetHandlerCases を、実際のケースで
// 呼び出す。name が空のときのバリデーションを忘れた実装（mutant.go）では
// 「name無し」のケースでステータスコードが不一致になり、ユーザーのテストが
// 正しく比較・失敗処理をしていれば検出できる。
func TestGreetHandler(t *testing.T) {
	cases := []greetCase{
		{"with name returns 200", "?name=Alice", http.StatusOK, "Hello, Alice!"},
		{"missing name returns 400", "", http.StatusBadRequest, ""},
		{"empty name value returns 400", "?name=", http.StatusBadRequest, ""},
	}
	runGreetHandlerCases(t, cases)
}
