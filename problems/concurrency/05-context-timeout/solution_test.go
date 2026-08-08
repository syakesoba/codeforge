package main

import (
	"testing"
	"time"
)

func TestSumWithDeadline(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	got, err := sumWithDeadline(nums, 500*time.Millisecond)

	if err != nil {
		t.Fatalf("エラーは発生しないはずですが、エラーが返りました: %v", err)
	}
	if got != 15 {
		t.Fatalf("sumWithDeadline(%v) = %d, want 15", nums, got)
	}
}

func TestSumWithDeadlineDifferentInput(t *testing.T) {
	nums := []int{10, -3, 7, 2}

	got, err := sumWithDeadline(nums, 500*time.Millisecond)

	if err != nil {
		t.Fatalf("エラーは発生しないはずですが、エラーが返りました: %v", err)
	}
	if got != 16 {
		t.Fatalf("sumWithDeadline(%v) = %d, want 16", nums, got)
	}
}
