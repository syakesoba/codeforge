package main

import (
	"reflect"
	"testing"
)

func TestSquareAll(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, -3, 0}
	want := []int{1, 4, 9, 16, 25, 9, 0}

	got := squareAll(input)

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("squareAll(%v) = %v, want %v", input, got, want)
	}
}

func TestSquareAllPreservesLength(t *testing.T) {
	input := make([]int, 50)
	for i := range input {
		input[i] = i
	}

	got := squareAll(input)

	if len(got) != len(input) {
		t.Fatalf("結果の長さが一致しません。期待値: %d, 実際: %d", len(input), len(got))
	}
	for i, n := range input {
		want := n * n
		if got[i] != want {
			t.Fatalf("squareAll[%d] = %d, want %d", i, got[i], want)
		}
	}
}
