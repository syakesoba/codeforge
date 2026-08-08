package main

import (
	"strings"
	"testing"
)

func TestHashPasswordProducesHash(t *testing.T) {
	hashed, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("HashPassword がエラーを返しました: %v", err)
	}
	if hashed == "" {
		t.Fatal("HashPassword が空文字列を返しました")
	}
	if hashed == "password123" {
		t.Fatal("パスワードがハッシュ化されていません（平文のまま返っています）")
	}
	// bcryptのハッシュは "$2a$" などのプレフィックスを持つ
	if !strings.HasPrefix(hashed, "$2") {
		t.Fatalf("bcryptのハッシュ形式ではありません。実際: %q", hashed)
	}
}

func TestHashPasswordIsSaltedEachTime(t *testing.T) {
	h1, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("HashPassword がエラーを返しました: %v", err)
	}
	h2, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("HashPassword がエラーを返しました: %v", err)
	}

	if h1 == h2 {
		t.Fatal("同じパスワードから同じハッシュが生成されています。bcrypt.GenerateFromPassword を使っていますか？")
	}
}

func TestCheckPasswordCorrect(t *testing.T) {
	hashed, err := HashPassword("correct-horse")
	if err != nil {
		t.Fatalf("HashPassword がエラーを返しました: %v", err)
	}

	if !CheckPassword(hashed, "correct-horse") {
		t.Fatal("正しいパスワードなのに false が返りました")
	}
}

func TestCheckPasswordWrong(t *testing.T) {
	hashed, err := HashPassword("correct-horse")
	if err != nil {
		t.Fatalf("HashPassword がエラーを返しました: %v", err)
	}

	if CheckPassword(hashed, "wrong-password") {
		t.Fatal("間違ったパスワードなのに true が返りました")
	}
}

func TestCheckPasswordWithInvalidHash(t *testing.T) {
	if CheckPassword("not-a-bcrypt-hash", "anything") {
		t.Fatal("不正なハッシュに対して true が返りました")
	}
}
