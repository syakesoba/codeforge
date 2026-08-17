//go:build ignore

package main

import (
	"errors"
	"fmt"
)

var ErrPermissionDenied = errors.New("permission denied")

type RangeError struct {
	Min, Max, Got int
}

func (e *RangeError) Error() string {
	return fmt.Sprintf("must be between %d and %d, got %d", e.Min, e.Max, e.Got)
}

func CheckAccess(role string, level int) error {
	if role != "admin" {
		return fmt.Errorf("check access: %w", ErrPermissionDenied)
	}
	if level < 1 || level > 10 {
		return fmt.Errorf("check access: %w", &RangeError{Min: 1, Max: 10, Got: level})
	}
	return nil
}

func Describe(err error) string {
	if err == nil {
		return "OK"
	}
	if errors.Is(err, ErrPermissionDenied) {
		return "権限がありません"
	}
	var rangeErr *RangeError
	if errors.As(err, &rangeErr) {
		return fmt.Sprintf("レベルは%d〜%dの範囲で指定してください（実際: %d）", rangeErr.Min, rangeErr.Max, rangeErr.Got)
	}
	return "不明なエラー"
}

func main() {}
