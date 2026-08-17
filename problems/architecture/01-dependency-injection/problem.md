# Lesson 1: インターフェースで依存を注入する

## 今回学ぶこと

- 「具体的な実装」ではなく「インターフェース」に依存するとは何か
- 依存性の注入（Dependency Injection, DI）の基本
- なぜこの設計がテストしやすさ・柔軟性につながるのか

### 具体的な実装に直接依存すると何が困るか

たとえば「注文を受け付けたらメールで通知する」処理を、素朴に書くとこうなります。

```go
type OrderService struct {
	emailClient *EmailClient // 具体的な実装に直接依存
}

func (s *OrderService) PlaceOrder(item string) error {
	return s.emailClient.SendEmail("ご注文ありがとうございます: " + item)
}
```

これだと、次のような問題が出てきます。

- 通知手段をメールからSlackに変えたくなったら、`OrderService` 自体を書き換える必要がある
- `PlaceOrder` をテストしたいだけなのに、本物のメール送信（外部への通信）が必要になってしまう

### インターフェースに依存する

そこで、「通知する」という**振る舞い**だけをインターフェースとして定義し、`OrderService` はそのインターフェースにだけ依存させます。

```go
type Notifier interface {
	Notify(message string) error
}

type OrderService struct {
	notifier Notifier // 具体的な実装ではなく、インターフェースに依存
}
```

`Notifier` インターフェースを満たしてさえいれば、`EmailNotifier` でも `SlackNotifier` でも、テスト用の偽物（フェイク）でも、何を渡しても `OrderService` は変更なしに動きます。

### コンストラクタで「注入」する

依存（この場合は `Notifier`）を構造体の外から渡すことを、**依存性の注入（DI）** と呼びます。Goでは、多くの場合コンストラクタ関数の引数として渡すだけのシンプルな形で実現します。

```go
func NewOrderService(notifier Notifier) *OrderService {
	return &OrderService{notifier: notifier}
}

// 呼び出し側が「どの実装を使うか」を決める
svc := NewOrderService(&EmailNotifier{})       // 本番ではメール通知
svc := NewOrderService(&fakeNotifierForTest{}) // テストでは偽物
```

`OrderService` 自身は、渡されたものが本物かテスト用の偽物かを一切気にしません。これが、Lesson 3で学ぶ「サービス層」を単体でテストしやすくする土台になります。

## 演習

`OrderService` と、それを組み立てる `NewOrderService`、注文処理を行う `PlaceOrder` を実装してください。

**`NewOrderService(notifier Notifier) *OrderService`**
- 受け取った `notifier` を保持した `*OrderService` を返す

**`(s *OrderService) PlaceOrder(item string) error`**
- `s.notifier.Notify(...)` を、`"ご注文ありがとうございます: " + item` に相当するメッセージで呼び出す
- `Notify` の戻り値（エラー）をそのまま返す

## ヒント

- `PlaceOrder` のメッセージ組み立てには `fmt.Sprintf("ご注文ありがとうございます: %s", item)` が使えます。`fmt` のimportを忘れずに（「インポートを自動修正」が使えます）
- `NewOrderService` は1行、`PlaceOrder` も実質1行（`return s.notifier.Notify(...)`）で書けます
