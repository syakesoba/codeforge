package main

import (
	"errors"
	"fmt"
)

// ErrPermissionDenied は権限不足を表すセンチネルエラーです。
var ErrPermissionDenied = errors.New("permission denied")

// RangeError は数値が許容範囲外だったことを表すエラー型です。
type RangeError struct {
	Min, Max, Got int
}

func (e *RangeError) Error() string {
	return fmt.Sprintf("must be between %d and %d, got %d", e.Min, e.Max, e.Got)
}

// CheckAccess はアクセス権限を検証します（実装済み・変更不要）。
// admin以外のroleならErrPermissionDeniedを、levelが範囲外ならRangeErrorを、
// どちらもラップして返します。
func CheckAccess(role string, level int) error {
	if role != "admin" {
		return fmt.Errorf("check access: %w", ErrPermissionDenied)
	}
	if level < 1 || level > 10 {
		return fmt.Errorf("check access: %w", &RangeError{Min: 1, Max: 10, Got: level})
	}
	return nil
}

// Describe は CheckAccess が返したエラーを、日本語のメッセージに変換します。
// err が nil なら "OK" を返してください。
func Describe(err error) string {
	// TODO:
	// 1. err が nil なら "OK" を返す
	// 2. errors.Is(err, ErrPermissionDenied) が true なら "権限がありません" を返す
	// 3. var rangeErr *RangeError を宣言し、errors.As(err, &rangeErr) が true なら
	//    fmt.Sprintf("レベルは%d〜%dの範囲で指定してください（実際: %d）", rangeErr.Min, rangeErr.Max, rangeErr.Got) を返す
	// 4. どれにも当てはまらなければ "不明なエラー" を返す
	return ""
}

func main() {}
