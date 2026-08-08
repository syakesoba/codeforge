package main

import (
	"reflect"
	"testing"
	"time"
)

func TestCollectSquares(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	want := []int{1, 4, 9, 16, 25}

	got := collectSquares(input)

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("collectSquares(%v) = %v, want %v", input, got, want)
	}
}

func TestCollectSquaresEmpty(t *testing.T) {
	got := collectSquares(nil)
	if len(got) != 0 {
		t.Fatalf("空スライスに対しては空の結果が期待されますが、%v が返りました", got)
	}
}

// collectSquares がチャネルをcloseし忘れて永遠にブロックしていないかを確認する。
func TestCollectSquaresDoesNotHang(t *testing.T) {
	done := make(chan struct{})
	go func() {
		collectSquares([]int{1, 2, 3})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("collectSquares が完了しませんでした。channelをcloseし忘れていませんか？")
	}
}
