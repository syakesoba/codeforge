package main

// Notifier は通知を送信する手段を表す抽象化です。
// 「メールで送る」「Slackに送る」といった具体的な実装の詳細を知らなくても、
// OrderServiceはこのインターフェースにだけ依存すればよくなります。
type Notifier interface {
	Notify(message string) error
}

// OrderService は注文処理のロジックを持ちます。
// Notifierを具体的な実装（例: EmailNotifier）に直接依存させるのではなく、
// インターフェースとして「外から受け取る」設計にします。
// これを依存性の注入（Dependency Injection, DI）と呼びます。
type OrderService struct {
	notifier Notifier
}

// NewOrderService は OrderService のコンストラクタです。
// 呼び出し側がどんなNotifier実装を渡すかを決められるようにします。
func NewOrderService(notifier Notifier) *OrderService {
	// TODO: &OrderService{notifier: notifier} を返す
	return nil
}

// PlaceOrder は注文を処理し、notifierを通じて通知します。
func (s *OrderService) PlaceOrder(item string) error {
	// TODO: s.notifier.Notify(fmt.Sprintf("ご注文ありがとうございます: %s", item)) を呼び、
	// その戻り値のエラーをそのまま返す
	return nil
}

func main() {}
