package main

import "testing"

func TestOrderStatusMessageSuccess(t *testing.T) {
	err := ProcessOrder(2, 10, 500, 2000)
	got := OrderStatusMessage(err)
	want := "注文が完了しました"
	if got != want {
		t.Errorf("OrderStatusMessage = %q, want %q", got, want)
	}
}

func TestOrderStatusMessageInvalidQuantity(t *testing.T) {
	err := ProcessOrder(0, 10, 500, 2000)
	got := OrderStatusMessage(err)
	want := "数量は1以上で指定してください"
	if got != want {
		t.Errorf("OrderStatusMessage = %q, want %q", got, want)
	}
}

func TestOrderStatusMessageNegativeQuantity(t *testing.T) {
	err := ProcessOrder(-3, 10, 500, 2000)
	got := OrderStatusMessage(err)
	want := "数量は1以上で指定してください"
	if got != want {
		t.Errorf("OrderStatusMessage = %q, want %q", got, want)
	}
}

func TestOrderStatusMessageOutOfStock(t *testing.T) {
	err := ProcessOrder(20, 10, 500, 100000)
	got := OrderStatusMessage(err)
	want := "在庫が不足しています"
	if got != want {
		t.Errorf("OrderStatusMessage = %q, want %q", got, want)
	}
}

func TestOrderStatusMessageInsufficientFunds(t *testing.T) {
	err := ProcessOrder(5, 10, 500, 1000)
	got := OrderStatusMessage(err)
	want := "残高が不足しています（必要: 2500円、残高: 1000円）"
	if got != want {
		t.Errorf("OrderStatusMessage = %q, want %q", got, want)
	}
}

func TestOrderStatusMessageExactBalance(t *testing.T) {
	err := ProcessOrder(2, 10, 500, 1000)
	got := OrderStatusMessage(err)
	want := "注文が完了しました"
	if got != want {
		t.Errorf("OrderStatusMessage = %q, want %q", got, want)
	}
}
