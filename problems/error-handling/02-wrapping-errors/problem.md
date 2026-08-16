# Lesson 2: エラーをラップする

## 今回学ぶこと

- なぜエラーに「コンテキスト」を付け加える必要があるのか
- `fmt.Errorf` と `%w` によるエラーのラップ
- ラップしても元のエラーを判定できる仕組み（さわり）

### エラーにコンテキストを付ける

深い呼び出し階層の下の方で発生したエラーを、そのまま呼び出し元まで伝播させると、「結局どこで・何が原因で失敗したのか」が分かりにくくなります。

```go
func FindUser(id int) (string, error) {
	name, ok := users[id]
	if !ok {
		return "", ErrNotFound // idの情報が失われる
	}
	return name, nil
}
```

これだと、呼び出し元が受け取るのは `"not found"` という情報だけです。「どの `id` で失敗したのか」が分かりません。

### fmt.Errorf でコンテキストを付ける

```go
func FindUser(id int) (string, error) {
	name, ok := users[id]
	if !ok {
		return "", fmt.Errorf("user %d: %w", id, ErrNotFound)
	}
	return name, nil
}
```

`fmt.Errorf` は `%s` と同じ感覚でエラーメッセージに値を埋め込めますが、**`%w`** という特別な動詞を使うと、埋め込んだエラー（ここでは `ErrNotFound`）を元のまま内部に保持した新しいエラーを作ります。

- `err.Error()` は `"user 99: not found"` のような、コンテキストを含む文字列になる
- それでいて、後述する `errors.Is(err, ErrNotFound)` は **true を返す**（`%w` で包んだ元のエラーを覚えているため）

`%s` や `%v` で埋め込んだ場合はこの「元のエラーを覚えている」効果が得られない点に注意してください。単なる文字列の埋め込みになってしまいます。

### errors.Is で判定する（次のレッスンで詳しく扱います）

```go
_, err := FindUser(99)
if errors.Is(err, ErrNotFound) {
	fmt.Println("ユーザーが見つかりませんでした")
}
```

`err` は `fmt.Errorf` で作られた新しいエラーですが、`%w` でラップされているため、`errors.Is` はラップの鎖をたどって `ErrNotFound` を見つけ出してくれます。この仕組みの詳細は次のレッスンで扱います。

## 演習

`FindUser` 関数を実装してください。

**`FindUser(id int) (string, error)`**
- `users` マップに `id` が存在すれば、対応する名前と `nil` を返す
- 存在しなければ、`fmt.Errorf("user %d: %w", id, ErrNotFound)` でラップしたエラーを返す

## ヒント

- マップからの検索は `name, ok := users[id]` の形（カンマOKイディオム）で行います
- `fmt.Errorf` の第2引数以降は `%w` の位置に対応する値を渡します。`id` は `%d`、`ErrNotFound` は `%w` です
