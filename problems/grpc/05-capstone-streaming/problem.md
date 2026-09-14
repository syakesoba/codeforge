# Lesson 5（道場）: ストリーミングRPCを実装する

## 今回学ぶこと

- サーバーストリーミングRPCとは何か
- 1回のリクエストに対して複数のレスポンスを順番に送る方法
- これまでのレッスンで学んだステータスコードによるエラー処理と組み合わせる

### サーバーストリーミングRPC

Lesson 1〜4で扱ってきたのは、1リクエストに1レスポンスが対応する**Unary RPC**でした。gRPCにはこれ以外に、1回のリクエストに対して**複数のレスポンスを順番に返す**方式があります。これを**サーバーストリーミングRPC**と呼びます。

`.proto` では、レスポンス側の型に `stream` キーワードを付けて宣言します。

```protobuf
service Counter {
  rpc CountUp(CountRequest) returns (stream CountReply);
}
```

大量データの分割送信や、進捗を逐次通知したい場合（ファイルの処理状況など）に向いた方式です。REST/JSONで同じことをやろうとすると、ポーリングやWebSocketなど別の仕組みが必要になりますが、gRPCでは同じ枠組みの中で自然に書けます。

### サーバー側の実装

サーバーストリーミングのハンドラは、通常の `(resp, error)` を返す形ではなく、`stream` 引数の `Send` メソッドを**複数回**呼び出す形になります。

```go
func (s *counterServer) CountUp(req *CountRequest, stream Counter_CountUpServer) error {
	for i := int32(1); i <= req.GetTo(); i++ {
		if err := stream.Send(&CountReply{Value: i}); err != nil {
			return err // 送信自体が失敗した場合（クライアント切断など）
		}
	}
	return nil // 戻り値がnilなら、クライアント側は最後にio.EOFを受け取る
}
```

クライアント側は `Recv()` を繰り返し呼び出し、`io.EOF` が返ってきたら受信完了です（Lesson 1・2で見た「1回呼んで終わり」のUnary RPCとの一番の違いです）。

## 演習

これまでのレッスンで学んだ「入力チェック＋ステータスコード」（Lesson 3）と「ストリーミング」を組み合わせて、`counterServer` の `CountUp` を実装してください。

**`(s *counterServer) CountUp(req *CountRequest, stream Counter_CountUpServer) error`**
1. `req.GetTo()` が `1` 未満の場合、`status.Error(codes.InvalidArgument, "to must be at least 1")` を返す
2. `1` から `req.GetTo()` まで1ずつ増やしながら、`stream.Send(&CountReply{Value: i})` を呼び出す。送信中にエラーが起きたら、そのエラーをそのまま返す
3. すべて送り終えたら `nil` を返す

## ヒント

- `google.golang.org/grpc/codes` と `google.golang.org/grpc/status` のインポートが必要です（「インポートを自動修正」が使えます）
- ループは `for i := int32(1); i <= req.GetTo(); i++ { ... }` の形になります（`CountRequest.To` / `CountReply.Value` は `int32`型です）
