//go:build ignore

package main

import (
	"errors"
	"testing"
)

// parseCase は parsePositiveInt の1テストケースを表します。
type parseCase struct {
	name      string
	input     string
	want      int
	wantErrIs error
}

// runParsePositiveIntCases は cases の各要素について parsePositiveInt を呼び出し、
// t.Run でサブテストとして実行し、値またはエラーが期待通りか検証してください。
func runParsePositiveIntCases(t *testing.T, cases []parseCase) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parsePositiveInt(c.input)

			if c.wantErrIs != nil {
				if !errors.Is(err, c.wantErrIs) {
					t.Errorf("parsePositiveInt(%q) error = %v, want error Is %v", c.input, err, c.wantErrIs)
				}
				return
			}

			if err != nil {
				t.Errorf("parsePositiveInt(%q) unexpected error: %v", c.input, err)
				return
			}
			if got != c.want {
				t.Errorf("parsePositiveInt(%q) = %d, want %d", c.input, got, c.want)
			}
		})
	}
}
