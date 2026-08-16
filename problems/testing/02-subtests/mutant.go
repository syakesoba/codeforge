//go:build ignore

package main

// daysInMonth はわざとバグを仕込んだ実装です
// （うるう年の判定を忘れており、2月は常に28日を返してしまう）。
func daysInMonth(month int, isLeapYear bool) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		return 28
	default:
		return 0
	}
}
