package main

import "testing"

func TestValidateAgeValid(t *testing.T) {
	if err := ValidateAge(30); err != nil {
		t.Fatalf("有効な年齢なのにエラーが返りました: %v", err)
	}
	if err := ValidateAge(0); err != nil {
		t.Fatalf("有効な年齢なのにエラーが返りました: %v", err)
	}
	if err := ValidateAge(150); err != nil {
		t.Fatalf("有効な年齢なのにエラーが返りました: %v", err)
	}
}

func TestValidateAgeNegative(t *testing.T) {
	err := ValidateAge(-1)
	if err == nil {
		t.Fatal("不正な年齢なのにエラーがnilでした")
	}
	if err.Error() != "age: 0から150の範囲で指定してください" {
		t.Fatalf("Error()の内容が想定と異なります。実際: %q", err.Error())
	}
}

func TestValidateAgeTooLarge(t *testing.T) {
	err := ValidateAge(151)
	if err == nil {
		t.Fatal("不正な年齢なのにエラーがnilでした")
	}
	if err.Error() != "age: 0から150の範囲で指定してください" {
		t.Fatalf("Error()の内容が想定と異なります。実際: %q", err.Error())
	}
}

func TestValidateAgeReturnsValidationErrorType(t *testing.T) {
	err := ValidateAge(-5)
	if err == nil {
		t.Fatal("不正な年齢なのにエラーがnilでした")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("戻り値が*ValidationErrorではありません。実際の型: %T", err)
	}
	if ve.Field != "age" {
		t.Errorf("Field = %q, want %q", ve.Field, "age")
	}
	if ve.Message == "" {
		t.Error("Messageが空です")
	}
}
