//go:build ignore

package main

import (
	"errors"
	"strconv"
)

var ErrNotPositive = errors.New("value must be positive")

// parsePositiveInt はわざとバグを仕込んだ実装です
// （0以下かどうかのチェックを忘れており、負の数や0もそのまま通してしまう）。
func parsePositiveInt(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return n, nil
}
