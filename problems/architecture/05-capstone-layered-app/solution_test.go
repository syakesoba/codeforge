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

func newTestLibrary() (*InMemoryBookRepository, *fakeNotifier, *LibraryService) {
	repo := NewInMemoryBookRepository(map[int]*Book{
		1: {ID: 1, Title: "Goならわかるシステムプログラミング", OnLoan: false},
		2: {ID: 2, Title: "リーダブルコード", OnLoan: true},
	})
	notifier := &fakeNotifier{}
	svc := NewLibraryService(repo, notifier)
	return repo, notifier, svc
}

func TestBorrowSuccess(t *testing.T) {
	repo, notifier, svc := newTestLibrary()

	if err := svc.Borrow(1); err != nil {
		t.Fatalf("Borrowがエラーを返しました: %v", err)
	}

	book, _ := repo.FindByID(1)
	if !book.OnLoan {
		t.Error("貸し出し後もOnLoanがfalseのままです")
	}

	if len(notifier.messages) != 1 {
		t.Fatalf("Notifyが呼ばれた回数 = %d, want 1", len(notifier.messages))
	}
	want := "「Goならわかるシステムプログラミング」を貸し出しました"
	if notifier.messages[0] != want {
		t.Errorf("通知メッセージ = %q, want %q", notifier.messages[0], want)
	}
}

func TestBorrowAlreadyOnLoan(t *testing.T) {
	repo, notifier, svc := newTestLibrary()

	err := svc.Borrow(2)
	if !errors.Is(err, ErrBookAlreadyOnLoan) {
		t.Fatalf("貸し出し中の本なのにErrBookAlreadyOnLoanが返りませんでした。実際: %v", err)
	}
	if len(notifier.messages) != 0 {
		t.Error("失敗したのに通知が送られています")
	}
	book, _ := repo.FindByID(2)
	if !book.OnLoan {
		t.Error("book.OnLoanの状態が変わってしまっています")
	}
}

func TestBorrowNotFound(t *testing.T) {
	_, notifier, svc := newTestLibrary()

	err := svc.Borrow(999)
	if !errors.Is(err, ErrBookNotFound) {
		t.Fatalf("存在しない本なのにErrBookNotFoundが返りませんでした。実際: %v", err)
	}
	if len(notifier.messages) != 0 {
		t.Error("失敗したのに通知が送られています")
	}
}

func TestBorrowPropagatesNotifierError(t *testing.T) {
	repo := NewInMemoryBookRepository(map[int]*Book{
		1: {ID: 1, Title: "テスト駆動開発", OnLoan: false},
	})
	notifyErr := errors.New("notification failed")
	notifier := &fakeNotifier{err: notifyErr}
	svc := NewLibraryService(repo, notifier)

	err := svc.Borrow(1)
	if !errors.Is(err, notifyErr) {
		t.Fatalf("notifierのエラーが伝播していません。実際: %v", err)
	}

	// 通知に失敗しても、貸し出し状態自体はすでに保存されているはず
	book, _ := repo.FindByID(1)
	if !book.OnLoan {
		t.Error("book.OnLoanがtrueになっていません")
	}
}
