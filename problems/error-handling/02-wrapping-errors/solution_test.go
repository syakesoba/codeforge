package main

import (
	"errors"
	"strings"
	"testing"
)

func TestFindUserFound(t *testing.T) {
	name, err := FindUser(1)
	if err != nil {
		t.Fatalf("存在するユーザーなのにエラーが返りました: %v", err)
	}
	if name != "Alice" {
		t.Errorf("name = %q, want %q", name, "Alice")
	}
}

func TestFindUserNotFound(t *testing.T) {
	_, err := FindUser(99)
	if err == nil {
		t.Fatal("存在しないユーザーなのにエラーがnilでした")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatal("errors.Is(err, ErrNotFound) がfalseです。%w でラップしていますか？")
	}
	if !strings.Contains(err.Error(), "99") {
		t.Errorf("エラーメッセージにidのコンテキストが含まれていません。実際: %q", err.Error())
	}
}

func TestFindUserAnotherNotFound(t *testing.T) {
	_, err := FindUser(3)
	if err == nil {
		t.Fatal("存在しないユーザーなのにエラーがnilでした")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatal("errors.Is(err, ErrNotFound) がfalseです。%w でラップしていますか？")
	}
}
