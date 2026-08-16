//go:build ignore

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// greetCase は greetHandler の1テストケースを表します。
type greetCase struct {
	name       string
	query      string
	wantStatus int
	wantBody   string
}

// runGreetHandlerCases は cases の各要素について、httptest でリクエストを作り
// greetHandler を呼び出し、ステータスコード・本文が期待通りか検証してください。
func runGreetHandlerCases(t *testing.T, cases []greetCase) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/greet"+c.query, nil)
			rec := httptest.NewRecorder()

			greetHandler(rec, req)

			if rec.Code != c.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, c.wantStatus)
			}
			if c.wantStatus == http.StatusOK && rec.Body.String() != c.wantBody {
				t.Errorf("body = %q, want %q", rec.Body.String(), c.wantBody)
			}
		})
	}
}
