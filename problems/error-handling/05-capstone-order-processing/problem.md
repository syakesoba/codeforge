# Lesson 5（道場）: 実践的なエラーハンドリング

## 今回学ぶこと

- これまでのレッスンで学んだ要素を組み合わせて、実践的なエラー処理を書く
- センチネルエラー・カスタムエラー型・ラップ・`errors.Is`/`errors.As` の使い分け

### これまで学んだことの総仕上げ

このレッスンは「道場」＝総仕上げ問題です。Lesson 1〜3で学んだ以下の要素を、すべて組み合わせて使います。

- カスタムエラー型（Lesson 1）: `InsufficientFundsError` が、必要金額・残高という追加情報を持つ
- エラーのラップ（Lesson 2）: センチネルエラーを `fmt.Errorf("...: %w", ...)` でラップしてコンテキストを付ける
- `errors.Is` / `errors.As`（Lesson 3）: エラーの種類によって、ユーザー向けメッセージを出し分ける

実際のアプリケーションでは、「関数の奥深くで発生したエラーを、最終的にユーザーが理解できるメッセージに変換する」という処理が頻繁に登場します。今回の `OrderStatusMessage` はまさにその役割です。

## 演習

以下の型・関数は**実装済み**です（変更不要）。

```go
var ErrInvalidQuantity = errors.New("invalid quantity")
var ErrOutOfStock = errors.New("out of stock")

type InsufficientFundsError struct {
	Needed    int
	Available int
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("insufficient funds: needed %d, available %d", e.Needed, e.Available)
}

// ProcessOrder は注文を処理します。
// quantityが1未満、stockが足りない、balanceが足りない、のいずれかなら
// それぞれ対応するエラーを返します。
func ProcessOrder(quantity, stock, price, balance int) error {
	if quantity < 1 {
		return fmt.Errorf("process order: %w", ErrInvalidQuantity)
	}
	if quantity > stock {
		return fmt.Errorf("process order: %w", ErrOutOfStock)
	}
	total := quantity * price
	if total > balance {
		return &InsufficientFundsError{Needed: total, Available: balance}
	}
	return nil
}
```

`OrderStatusMessage(err error) string` を実装してください。`ProcessOrder` が返す4種類の結果（成功・数量不正・在庫不足・残高不足）それぞれに対して、以下のメッセージを返します。

| ケース | 返すメッセージ |
|---|---|
| `err == nil` | `"注文が完了しました"` |
| `errors.Is(err, ErrInvalidQuantity)` | `"数量は1以上で指定してください"` |
| `errors.Is(err, ErrOutOfStock)` | `"在庫が不足しています"` |
| `errors.As` で `*InsufficientFundsError` を取得 | `fmt.Sprintf("残高が不足しています（必要: %d円、残高: %d円）", fundsErr.Needed, fundsErr.Available)` |
| それ以外 | `"予期しないエラーが発生しました"` |

## ヒント

- 判定の順序は上の表の通りで問題ありません。`err == nil` を最初に確認してください
- `errors.As` の第2引数は `&fundsErr`（`*InsufficientFundsError` 型変数へのポインタ）です
- `InsufficientFundsError` は `errors.New` で作られたセンチネルエラーではなく独自の型なので `errors.Is` ではなく `errors.As` で判定します
