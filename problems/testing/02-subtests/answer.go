//go:build ignore

package main

import (
	"fmt"
	"testing"
)

// daysCase は daysInMonth の1テストケースを表します。
type daysCase struct {
	month      int
	isLeapYear bool
	want       int
}

// runDaysInMonthCases は cases の各要素について daysInMonth を呼び出し、
// t.Run でサブテストとして実行し、結果が want と一致するか検証してください。
func runDaysInMonthCases(t *testing.T, cases []daysCase) {
	for _, c := range cases {
		name := fmt.Sprintf("month=%d/leap=%v", c.month, c.isLeapYear)
		t.Run(name, func(t *testing.T) {
			got := daysInMonth(c.month, c.isLeapYear)
			if got != c.want {
				t.Errorf("daysInMonth(%d, %v) = %d, want %d", c.month, c.isLeapYear, got, c.want)
			}
		})
	}
}
