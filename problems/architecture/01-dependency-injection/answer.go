//go:build ignore

package main

import "fmt"

type Notifier interface {
	Notify(message string) error
}

type OrderService struct {
	notifier Notifier
}

func NewOrderService(notifier Notifier) *OrderService {
	return &OrderService{notifier: notifier}
}

func (s *OrderService) PlaceOrder(item string) error {
	return s.notifier.Notify(fmt.Sprintf("ご注文ありがとうございます: %s", item))
}

func main() {}
