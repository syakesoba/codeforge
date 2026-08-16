package main

import (
	"errors"
	"testing"
)

type fakeNotifier struct {
	messages []string
	err      error
}

func (f *fakeNotifier) Notify(message string) error {
	f.messages = append(f.messages, message)
	return f.err
}

func TestPlaceOrderNotifiesWithMessage(t *testing.T) {
	notifier := &fakeNotifier{}
	svc := NewOrderService(notifier)

	if err := svc.PlaceOrder("ノートPC"); err != nil {
		t.Fatalf("PlaceOrderがエラーを返しました: %v", err)
	}

	if len(notifier.messages) != 1 {
		t.Fatalf("Notifyが呼ばれた回数 = %d, want 1", len(notifier.messages))
	}
	want := "ご注文ありがとうございます: ノートPC"
	if notifier.messages[0] != want {
		t.Errorf("Notifyに渡されたメッセージ = %q, want %q", notifier.messages[0], want)
	}
}

func TestPlaceOrderPropagatesNotifierError(t *testing.T) {
	wantErr := errors.New("network error")
	notifier := &fakeNotifier{err: wantErr}
	svc := NewOrderService(notifier)

	err := svc.PlaceOrder("マウス")
	if !errors.Is(err, wantErr) {
		t.Fatalf("PlaceOrderがnotifierのエラーを伝播していません。実際: %v", err)
	}
}

func TestNewOrderServiceStoresNotifier(t *testing.T) {
	notifier := &fakeNotifier{}
	svc := NewOrderService(notifier)
	if svc == nil {
		t.Fatal("NewOrderServiceがnilを返しました")
	}
	// PlaceOrder経由で実際にnotifierが使われているかを確認する
	_ = svc.PlaceOrder("キーボード")
	if len(notifier.messages) != 1 {
		t.Fatal("コンストラクタに渡したnotifierが使われていません")
	}
}
