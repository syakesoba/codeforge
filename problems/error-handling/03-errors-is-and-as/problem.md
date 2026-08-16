# Lesson 3: errors.Is / errors.As で判定する

## 今回学ぶこと

- `errors.Is`: ラップの鎖をたどってセンチネルエラーを判定する
- `errors.As`: ラップの鎖をたどって特定の型のエラーを取り出す
- なぜ `==` や型アサーションを直接使ってはいけないのか

### なぜ err == ErrXxx が使えないのか

Lesson 2で見たように、`fmt.Errorf("...: %w", ErrNotFound)` は `ErrNotFound` そのものではなく、それを内部に包んだ**別の**エラー値を作ります。そのため、単純な比較では判定できません。

```go
err := fmt.Errorf("check access: %w", ErrPermissionDenied)
err == ErrPermissionDenied // false！ errはfmt.Errorfが作った別の値
```

同様に、`err.(*RangeError)` のような型アサーションも、`err` がラップされたエラーだと失敗します（`err` の実際の型は `*fmt.wrapError` であって `*RangeError` ではないため）。

### errors.Is: センチネルエラーの判定

`errors.Is(err, target)` は、`err` から始めて `%w` で包まれた鎖を1つずつたどりながら、`target` と一致するものが見つかるまで探索します。

```go
err := CheckAccess("guest", 5) // 内部で fmt.Errorf("check access: %w", ErrPermissionDenied)
if errors.Is(err, ErrPermissionDenied) {
	fmt.Println("権限がありません")
}
```

`err` 自体は `ErrPermissionDenied` と等しくありませんが、その内部（`Unwrap()`）をたどると `ErrPermissionDenied` にたどり着くため、`errors.Is` は `true` を返します。

### errors.As: 特定の型のエラーを取り出す

一方、「エラーの種類」ではなく「エラーが持つ追加情報」を取り出したいときは `errors.As` を使います。

```go
var rangeErr *RangeError
if errors.As(err, &rangeErr) {
	fmt.Println("範囲外:", rangeErr.Min, rangeErr.Max, rangeErr.Got)
}
```

`errors.As` も同様にラップの鎖をたどりますが、`target`（第2引数、ポインタで渡す）と**同じ型**のエラーを見つけたら、そのエラー値を `target` に代入し `true` を返します。見つからなければ `false` を返し、`rangeErr` は `nil` のままです。

> **使い分け**: 「特定の1つのエラー値かどうか」を知りたいときは `errors.Is`、「特定の型のエラーで、その中身（フィールド）が欲しい」ときは `errors.As` を使います。

## 演習

`CheckAccess` は実装済みです（`role` が `"admin"` でなければ `ErrPermissionDenied` を、`level` が1〜10の範囲外なら `*RangeError` を、それぞれラップして返します）。

`Describe(err error) string` を実装してください。

- `err` が `nil` なら `"OK"` を返す
- `errors.Is(err, ErrPermissionDenied)` が `true` なら `"権限がありません"` を返す
- `errors.As` で `*RangeError` を取り出せたら、`fmt.Sprintf("レベルは%d〜%dの範囲で指定してください（実際: %d）", rangeErr.Min, rangeErr.Max, rangeErr.Got)` を返す
- どれにも当てはまらなければ `"不明なエラー"` を返す

## ヒント

- `errors.As` の第2引数は「取り出したい型のポインタのポインタ」です。`var rangeErr *RangeError` と宣言してから `errors.As(err, &rangeErr)` のように渡します
- 判定の順序が重要です。先に `err == nil` をチェックしないと、後続の `errors.Is`/`errors.As` が `nil` を渡されてパニックすることはありませんが、意図通りの分岐にならなくなります
