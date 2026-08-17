# Lesson 4: 複数のエラーをまとめる（errors.Join）

## 今回学ぶこと

- 複数の検証エラーを1つにまとめて返す `errors.Join`（Go 1.20〜）
- `errors.Join` で作ったエラーに対しても `errors.Is`/`errors.As` が機能すること

### 「最初の1件で打ち切る」か「まとめて返す」か

これまでのレッスンでは、検証に失敗したら即座に1つのエラーを返す（early return）パターンを扱いました。しかしフォーム入力の検証のように「複数の項目を一度にチェックして、問題があるものを全部まとめてユーザーに教えたい」場面もあります。

1件ずつ `if err != nil { return err }` で早期returnしてしまうと、ユーザーは1つのエラーを直してから再送信し、次のエラーに気づく…という繰り返しになり体験が悪くなります。

### errors.Join でまとめる

```go
var errs []error
if name == "" {
	errs = append(errs, errors.New("name is required"))
}
if age < 0 || age > 150 {
	errs = append(errs, errors.New("age out of range"))
}

return errors.Join(errs...)
```

`errors.Join(errs...)` は、渡された複数のエラーを1つの `error` にまとめます。

- `errs` が空、または全要素が `nil` なら `errors.Join` は **`nil`** を返します。だから「問題が0件なら自動的に成功扱いになる」という都合の良い性質があり、`return errors.Join(errs...)` の1行で「問題があればエラー、無ければnil」を両方表現できます
- まとめられたエラーの `Error()` は、各エラーメッセージを改行区切りで連結した文字列になります

### errors.Is / errors.As はJoinの中もたどれる

Go 1.20以降、`errors.Is` と `errors.As` は `errors.Join` で作られたエラーの**中身も再帰的に**探索するように拡張されています。つまり、まとめられた複数のエラーのうちどれか1つでもセンチネルエラーと一致すれば `errors.Is` は `true` を返します。この性質のおかげで、「まとめて返す」設計に変えても、呼び出し側の判定コードを大きく変える必要はありません。

## 演習

`ValidateForm(name string, age int, email string) error` を実装してください。

- `name` が空文字列なら `errors.New("name is required")` を追加する
- `age` が `0`未満または`150`より大きいなら `errors.New("age out of range")` を追加する
- `email` に `"@"` が含まれていなければ `errors.New("invalid email")` を追加する（`strings.Contains` を使う）
- 集めたエラーを `errors.Join` でまとめて返す（問題が無ければ自動的に `nil` になります）

## ヒント

- `var errs []error` と `errs = append(errs, ...)` の組み合わせで、条件に応じてエラーを集めます
- 最後は `return errors.Join(errs...)` の1行だけです。`errs` が空でも安全に呼べます
- `errors` と `strings` の2つのimportが必要になるので、書き終えたら「インポートを自動修正」を使ってください
