package main

import (
	"errors"
	"testing"
)

func newTestAccounts() *InMemoryAccountRepository {
	return NewInMemoryAccountRepository(map[int]*Account{
		1: {ID: 1, Balance: 1000},
		2: {ID: 2, Balance: 500},
	})
}

func TestTransferSuccess(t *testing.T) {
	repo := newTestAccounts()
	svc := NewTransferService(repo)

	if err := svc.Transfer(1, 2, 300); err != nil {
		t.Fatalf("Transferがエラーを返しました: %v", err)
	}

	from, _ := repo.FindByID(1)
	to, _ := repo.FindByID(2)
	if from.Balance != 700 {
		t.Errorf("送金元残高 = %d, want 700", from.Balance)
	}
	if to.Balance != 800 {
		t.Errorf("送金先残高 = %d, want 800", to.Balance)
	}
}

func TestTransferInsufficientBalance(t *testing.T) {
	repo := newTestAccounts()
	svc := NewTransferService(repo)

	err := svc.Transfer(2, 1, 10000)
	if !errors.Is(err, ErrInsufficientBalance) {
		t.Fatalf("残高不足なのにErrInsufficientBalanceが返りませんでした。実際: %v", err)
	}

	// 残高は変更されていないはず
	from, _ := repo.FindByID(2)
	to, _ := repo.FindByID(1)
	if from.Balance != 500 {
		t.Errorf("送金失敗時に送金元残高が変更されています。実際: %d", from.Balance)
	}
	if to.Balance != 1000 {
		t.Errorf("送金失敗時に送金先残高が変更されています。実際: %d", to.Balance)
	}
}

func TestTransferFromAccountNotFound(t *testing.T) {
	repo := newTestAccounts()
	svc := NewTransferService(repo)

	err := svc.Transfer(999, 1, 100)
	if err == nil {
		t.Fatal("存在しない送金元なのにエラーがnilでした")
	}
}

func TestTransferToAccountNotFound(t *testing.T) {
	repo := newTestAccounts()
	svc := NewTransferService(repo)

	err := svc.Transfer(1, 999, 100)
	if err == nil {
		t.Fatal("存在しない送金先なのにエラーがnilでした")
	}
	// 送金元の残高も変更されていないはず（送金先確認前に副作用が起きてはいけない）
	from, _ := repo.FindByID(1)
	if from.Balance != 1000 {
		t.Errorf("送金失敗時に送金元残高が変更されています。実際: %d", from.Balance)
	}
}

func TestTransferExactBalance(t *testing.T) {
	repo := newTestAccounts()
	svc := NewTransferService(repo)

	if err := svc.Transfer(2, 1, 500); err != nil {
		t.Fatalf("残高ちょうどの送金でエラーが返りました: %v", err)
	}
	from, _ := repo.FindByID(2)
	if from.Balance != 0 {
		t.Errorf("送金元残高 = %d, want 0", from.Balance)
	}
}
