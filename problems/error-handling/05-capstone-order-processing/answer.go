//go:build ignore

package main

import (
	"errors"
	"fmt"
)

var ErrInvalidQuantity = errors.New("invalid quantity")
var ErrOutOfStock = errors.New("out of stock")

type InsufficientFundsError struct {
	Needed    int
	Available int
}

func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("insufficient funds: needed %d, available %d", e.Needed, e.Available)
}

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

func OrderStatusMessage(err error) string {
	if err == nil {
		return "注文が完了しました"
	}
	if errors.Is(err, ErrInvalidQuantity) {
		return "数量は1以上で指定してください"
	}
	if errors.Is(err, ErrOutOfStock) {
		return "在庫が不足しています"
	}
	var fundsErr *InsufficientFundsError
	if errors.As(err, &fundsErr) {
		return fmt.Sprintf("残高が不足しています（必要: %d円、残高: %d円）", fundsErr.Needed, fundsErr.Available)
	}
	return "予期しないエラーが発生しました"
}

func main() {}
