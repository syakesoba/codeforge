package main

import "errors"

// ErrNotFound はユーザーが見つからなかったことを表すセンチネルエラーです。
// 呼び出し側は errors.Is(err, ErrNotFound) でこのエラーかどうかを判定できます。
var ErrNotFound = errors.New("not found")

var users = map[int]string{
	1: "Alice",
	2: "Bob",
}

// FindUser は id に対応するユーザー名を返します。
// 見つからない場合は、"user %d: " のようなコンテキストを付けつつ
// ErrNotFound を %w でラップしたエラーを返してください
// （ラップすることで、呼び出し側は errors.Is(err, ErrNotFound) で判定できるようになります）。
func FindUser(id int) (string, error) {
	// TODO:
	// name, ok := users[id] で検索する
	// 見つかれば name, nil を返す
	// 見つからなければ "", fmt.Errorf("user %d: %w", id, ErrNotFound) を返す
	return "", nil
}

func main() {}
