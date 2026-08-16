package main

import (
	"errors"
	"testing"
)

func TestFindByIDFound(t *testing.T) {
	repo := NewInMemoryUserRepository(map[int]*User{
		1: {ID: 1, Name: "Alice"},
	})

	user, err := repo.FindByID(1)
	if err != nil {
		t.Fatalf("存在するユーザーなのにエラーが返りました: %v", err)
	}
	if user.Name != "Alice" {
		t.Errorf("user.Name = %q, want %q", user.Name, "Alice")
	}
}

func TestFindByIDNotFound(t *testing.T) {
	repo := NewInMemoryUserRepository(map[int]*User{})

	_, err := repo.FindByID(99)
	if err == nil {
		t.Fatal("存在しないユーザーなのにエラーがnilでした")
	}
	if !errors.Is(err, ErrUserNotFound) {
		t.Error("errors.Is(err, ErrUserNotFound) がfalseです。%w でラップしていますか？")
	}
}

func TestUserNamesSuccess(t *testing.T) {
	repo := NewInMemoryUserRepository(map[int]*User{
		1: {ID: 1, Name: "Alice"},
		2: {ID: 2, Name: "Bob"},
		3: {ID: 3, Name: "Carol"},
	})

	names, err := UserNames(repo, []int{1, 3, 2})
	if err != nil {
		t.Fatalf("UserNamesがエラーを返しました: %v", err)
	}
	want := []string{"Alice", "Carol", "Bob"}
	if len(names) != len(want) {
		t.Fatalf("names = %v, want %v", names, want)
	}
	for i := range want {
		if names[i] != want[i] {
			t.Errorf("names[%d] = %q, want %q", i, names[i], want[i])
		}
	}
}

func TestUserNamesPropagatesError(t *testing.T) {
	repo := NewInMemoryUserRepository(map[int]*User{
		1: {ID: 1, Name: "Alice"},
	})

	_, err := UserNames(repo, []int{1, 999})
	if err == nil {
		t.Fatal("存在しないidが含まれるのにエラーがnilでした")
	}
	if !errors.Is(err, ErrUserNotFound) {
		t.Error("errors.Is(err, ErrUserNotFound) がfalseです")
	}
}

// fakeUserRepository は UserRepository インターフェースを満たす、テスト専用の別実装です。
// UserNames が InMemoryUserRepository という具体型ではなく UserRepository インターフェースに
// 依存していれば、この別実装を渡しても問題なく動作するはずです。
type fakeUserRepository struct {
	name string
}

func (f *fakeUserRepository) FindByID(id int) (*User, error) {
	return &User{ID: id, Name: f.name}, nil
}

func TestUserNamesAcceptsAnyRepositoryImplementation(t *testing.T) {
	repo := &fakeUserRepository{name: "Fake"}
	names, err := UserNames(repo, []int{1, 2})
	if err != nil {
		t.Fatalf("UserNamesがエラーを返しました: %v", err)
	}
	if len(names) != 2 || names[0] != "Fake" || names[1] != "Fake" {
		t.Errorf("names = %v, want [Fake Fake]", names)
	}
}
