package main

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateFormAllValid(t *testing.T) {
	if err := ValidateForm("Alice", 30, "alice@example.com"); err != nil {
		t.Fatalf("有効な入力なのにエラーが返りました: %v", err)
	}
}

func TestValidateFormOneInvalid(t *testing.T) {
	err := ValidateForm("", 30, "alice@example.com")
	if err == nil {
		t.Fatal("nameが空なのにエラーがnilでした")
	}
	if !strings.Contains(err.Error(), "name is required") {
		t.Errorf("エラーメッセージに'name is required'が含まれていません。実際: %q", err.Error())
	}
}

func TestValidateFormMultipleInvalid(t *testing.T) {
	err := ValidateForm("", 200, "not-an-email")
	if err == nil {
		t.Fatal("複数の問題があるのにエラーがnilでした")
	}
	msg := err.Error()
	for _, want := range []string{"name is required", "age out of range", "invalid email"} {
		if !strings.Contains(msg, want) {
			t.Errorf("エラーメッセージに %q が含まれていません。実際: %q", want, msg)
		}
	}
}

func TestValidateFormErrorsIsWorksThroughJoin(t *testing.T) {
	sentinel := errors.New("marker")
	joined := errors.Join(errors.New("other"), sentinel)
	if !errors.Is(joined, sentinel) {
		t.Fatal("errors.Join で作ったエラーに対して errors.Is が機能していません（テスト環境側の前提確認用）")
	}

	err := ValidateForm("", -5, "bad")
	if err == nil {
		t.Fatal("複数の問題があるのにエラーがnilでした")
	}
}

func TestValidateFormAgeOnly(t *testing.T) {
	err := ValidateForm("Bob", -1, "bob@example.com")
	if err == nil {
		t.Fatal("ageが範囲外なのにエラーがnilでした")
	}
	if !strings.Contains(err.Error(), "age out of range") {
		t.Errorf("エラーメッセージに'age out of range'が含まれていません。実際: %q", err.Error())
	}
	if strings.Contains(err.Error(), "name is required") {
		t.Error("nameは有効なのに'name is required'が含まれています")
	}
}
