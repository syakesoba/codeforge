# Lesson 2: gRPCクライアントを実装する

## 今回学ぶこと

- 生成されたクライアントスタブ（`GreeterClient`）の使い方
- gRPCの呼び出しはGoの関数呼び出しに近い形で書けること
- gRPCのエラーを、これまでのコースで学んだ `%w` によるラップと組み合わせて扱う方法

### クライアントスタブを呼び出す

Lesson 1では `Greeter` サービスの**サーバー側**を実装しました。今回は**クライアント側**です。`.proto` から生成された `GreeterClient` インターフェースには、サービス定義そのままの形でメソッドが並んでいます。

```go
type GreeterClient interface {
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}
```

呼び出し側は、通信の中身（HTTP/2やprotobufのシリアライズ）を一切意識せず、まるでローカルの関数を呼ぶかのように書けます。

```go
resp, err := client.SayHello(ctx, &HelloRequest{Name: "World"})
```

これがgRPCの大きな利点の1つです。REST/JSONのAPIをクライアントから呼ぶ場合、リクエストのJSONを組み立て、`http.Post` し、レスポンスのJSONをパースし……という手順が必要ですが、gRPCではその部分がすべて生成コードに隠蔽されています。

### gRPCのエラーをラップする

`client.SayHello` が返す `error` は、サーバー側で `status.Error(codes.XXX, "...")` のように作られたgRPC特有のエラー（ステータスコード付き）です（ステータスコードによるエラーの返し方はLesson 3で学びます）。

これまでの「エラーハンドリング」コースで学んだ `fmt.Errorf("...: %w", err)` によるラップは、gRPCのエラーに対しても同じように使えます。ラップしても、`status.Code(err)` で元のステータスコードを取り出せる仕組みになっています（`%w` でラップされたエラーチェーンを `status` パッケージがたどってくれるためです）。

```go
resp, err := client.SayHello(ctx, req)
if err != nil {
	return "", fmt.Errorf("call SayHello: %w", err)
}
```

呼び出し元の文脈（「どの呼び出しで失敗したか」）をエラーメッセージに残しつつ、ステータスコードによる判定はそのまま可能な状態を保てます。

## 演習

`CallSayHello` 関数を実装してください。

**`CallSayHello(ctx context.Context, client GreeterClient, name string) (string, error)`**
- `client.SayHello(ctx, &HelloRequest{Name: name})` を呼び出す
- エラーが返ってきたら、`"call SayHello: "` を前置して `%w` でラップして返す
- 成功したら、レスポンスのメッセージ文字列（`resp.GetMessage()`）と `nil` を返す

## ヒント

- `fmt.Errorf("call SayHello: %w", err)` の形でラップします。`fmt` のimportを忘れずに（「インポートを自動修正」が使えます）
- 成功時は `return resp.GetMessage(), nil`、失敗時は `return "", fmt.Errorf(...)` という形になります
