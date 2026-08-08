package main

import "testing"

func TestFirstResultCh1First(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 100

	got := firstResult(ch1, ch2)

	if got != 100 {
		t.Fatalf("firstResult() = %d, want 100", got)
	}
}

func TestFirstResultCh2First(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch2 <- 200

	got := firstResult(ch1, ch2)

	if got != 200 {
		t.Fatalf("firstResult() = %d, want 200", got)
	}
}
