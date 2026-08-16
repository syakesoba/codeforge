package main

import "testing"

func TestDescribeNil(t *testing.T) {
	if got := Describe(nil); got != "OK" {
		t.Errorf("Describe(nil) = %q, want %q", got, "OK")
	}
}

func TestDescribePermissionDenied(t *testing.T) {
	err := CheckAccess("guest", 5)
	got := Describe(err)
	want := "権限がありません"
	if got != want {
		t.Errorf("Describe(err) = %q, want %q", got, want)
	}
}

func TestDescribeRangeError(t *testing.T) {
	err := CheckAccess("admin", 99)
	got := Describe(err)
	want := "レベルは1〜10の範囲で指定してください（実際: 99）"
	if got != want {
		t.Errorf("Describe(err) = %q, want %q", got, want)
	}
}

func TestDescribeRangeErrorTooLow(t *testing.T) {
	err := CheckAccess("admin", 0)
	got := Describe(err)
	want := "レベルは1〜10の範囲で指定してください（実際: 0）"
	if got != want {
		t.Errorf("Describe(err) = %q, want %q", got, want)
	}
}

func TestDescribeValid(t *testing.T) {
	err := CheckAccess("admin", 5)
	if got := Describe(err); got != "OK" {
		t.Errorf("Describe(err) = %q, want %q", got, "OK")
	}
}
