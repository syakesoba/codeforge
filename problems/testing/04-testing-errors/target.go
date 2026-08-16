package main

import (
	"errors"
	"strconv"
)

// ErrNotPositive は、パース自体は成功したが値が正の整数ではなかった場合に
// 返されるエラーです（実装済み・変更不要）。
var ErrNotPositive = errors.New("value must be positive")

// parsePositiveInt は文字列を正の整数としてパースします（実装済み・変更不要）。
// パース自体に失敗した場合は strconv.Atoi のエラーを、
// パースできても0以下だった場合は ErrNotPositive を返します。
func parsePositiveInt(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if n <= 0 {
		return 0, ErrNotPositive
	}
	return n, nil
}
