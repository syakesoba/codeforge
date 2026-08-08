package main

import (
	"sort"
	"testing"
	"time"
)

func TestProcessJobs(t *testing.T) {
	jobs := []int{1, 2, 3, 4, 5}
	square := func(n int) int { return n * n }

	got := processJobs(jobs, 3, square)
	sort.Ints(got)

	want := []int{1, 4, 9, 16, 25}
	if len(got) != len(want) {
		t.Fatalf("結果の件数が一致しません。期待値: %v, 実際: %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("processJobs(%v) = %v, want %v (順不同で比較)", jobs, got, want)
		}
	}
}

func TestProcessJobsMoreWorkersThanJobs(t *testing.T) {
	jobs := []int{10, 20}
	double := func(n int) int { return n * 2 }

	got := processJobs(jobs, 8, double)
	sort.Ints(got)

	want := []int{20, 40}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("processJobs(%v) = %v, want %v (順不同で比較)", jobs, got, want)
		}
	}
}

// jobsCh / resultsCh のcloseが正しく連鎖しているかを確認する
// （closeし忘れているとrangeが永遠にブロックしてハングする）。
func TestProcessJobsDoesNotHang(t *testing.T) {
	done := make(chan struct{})
	go func() {
		processJobs([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 4, func(n int) int { return n })
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("processJobs が完了しませんでした。channelのcloseし忘れがないか確認してください。")
	}
}
