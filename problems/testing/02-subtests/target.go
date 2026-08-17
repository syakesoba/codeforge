package main

// daysInMonth は指定した月（1〜12）の日数を返します（実装済み・変更不要）。
// うるう年（isLeapYear）の場合、2月は29日になります。
func daysInMonth(month int, isLeapYear bool) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		if isLeapYear {
			return 29
		}
		return 28
	default:
		return 0
	}
}
