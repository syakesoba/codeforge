package main

import (
	"errors"
	"fmt"
)

// ErrInvalidQuantity は数量が不正だったことを表すセンチネルエラーです。
var ErrInvalidQuantity = errors.New("invalid quantity")

// ErrOutOfStock は在庫が足りないことを表すセンチネルエラーです。
var ErrOutOfStock = errors.New("out of stock")

// InsufficientFundsError は残高不足を表すエラー型です。
// 必要な金額と実際の残高の両方を保持します。
type InsufficientFundsError struct {
	Needed    int
	Available int
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("insufficient funds: needed %d, available %d", e.Needed, e.Available)
}

// ProcessOrder は注文を処理します（実装済み・変更不要）。
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

// OrderStatusMessage は ProcessOrder が返したエラーを、ユーザー向けの
// 日本語メッセージに変換してください。
func OrderStatusMessage(err error) string {
	// TODO:
	// 1. err が nil なら "注文が完了しました" を返す
	// 2. errors.Is(err, ErrInvalidQuantity) なら "数量は1以上で指定してください" を返す
	// 3. errors.Is(err, ErrOutOfStock) なら "在庫が不足しています" を返す
	// 4. var fundsErr *InsufficientFundsError を宣言し、errors.As(err, &fundsErr) が true なら
	//    fmt.Sprintf("残高が不足しています（必要: %d円、残高: %d円）", fundsErr.Needed, fundsErr.Available) を返す
	// 5. どれにも当てはまらなければ "予期しないエラーが発生しました" を返す
	return ""
}

func main() {}
