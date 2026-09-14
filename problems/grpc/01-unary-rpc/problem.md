# Lesson 1: Unary RPCサーバーを実装する

## 今回学ぶこと

- gRPCとは何か、REST/JSONのAPIと何が違うのか
- Protocol Buffers（`.proto`）でサービスとメッセージを定義し、コードを自動生成する仕組み
- 最も基本的な呼び出し方式「Unary RPC」のサーバー実装

### gRPCとは

これまでのコースでは `net/http` や Gin を使い、JSONをやり取りするWeb APIを作ってきました。gRPCはそれとは別系統の、サービス間通信のためのRPC（Remote Procedure Call）フレームワークです。

- **通信形式**: JSONではなく、Protocol Buffers（protobuf）というバイナリ形式でやり取りする。データが小さく、エンコード/デコードも速い
- **契約（インターフェース）を先に決める**: `.proto` ファイルにサービスの入出力の型を定義し、そこからサーバー用・クライアント用のコードを**自動生成**する。手書きのJSONパース処理や、型の取り違えによるバグが起きにくい
- **マイクロサービス間通信で特によく使われる**: 社内の複数サービス同士が通信するような場面で、REST/JSONより高速・型安全という理由で採用されることが多い

### `.proto` でサービスを定義する

このレッスンでは、以下のような `.proto` ファイルから、サーバー・クライアント両方のGoコードが自動生成されています（生成されたコードはこのレッスンのワークスペースに既に用意されており、あなたが書く必要はありません）。

```protobuf
service Greeter {
  rpc SayHello(HelloRequest) returns (HelloReply);
}

message HelloRequest {
  string name = 1;
}

message HelloReply {
  string message = 1;
}
```

- `service Greeter` が、提供するRPCの一覧（ここでは `SayHello` のみ）を宣言する
- `rpc SayHello(HelloRequest) returns (HelloReply);` は「`HelloRequest` を受け取り `HelloReply` を返す」という**契約**
- `message` はやり取りするデータの型。protoc（protobufのコンパイラ）がこれを読み、Goの構造体・シリアライズ処理・gRPC通信処理を自動生成する

生成されたコードには、次のようなものが含まれています。

- `HelloRequest` / `HelloReply`: `.proto` の `message` に対応するGo構造体（`GetName()` のようなgetterも自動で付く）
- `GreeterServer`: サーバーが実装すべきインターフェース（`SayHello(context.Context, *HelloRequest) (*HelloReply, error)`）
- `UnimplementedGreeterServer`: `GreeterServer` の全メソッドを「未実装」として満たす空実装。これを埋め込んでおくと、将来サービス定義にメソッドが増えても、実装し忘れているメソッドがあるだけではコンパイルエラーにならない（gRPCサーバー実装の定石）
- `RegisterGreeterServer`: 実装したサーバーを `*grpc.Server` に登録する関数

### Unary RPC

`SayHello` のように「1つのリクエストを送り、1つのレスポンスを受け取って終わり」という最もシンプルな呼び出し方式を **Unary RPC** と呼びます（普段のHTTP APIの1リクエスト1レスポンスに近いイメージです）。Lesson 5では、複数のレスポンスを順に受け取る「ストリーミングRPC」も扱います。

## 演習

`Greeter` サービスの実装である `server` 型の `SayHello` メソッドを実装してください。

**`(s *server) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error)`**
- `req.GetName()` で受け取った名前を使い、`"Hello, " + 名前 + "!"` という文字列を組み立てる
- そのメッセージを持つ `&HelloReply{Message: ...}` と `nil`（エラーなし）を返す

## ヒント

- `req.GetName()` は、`req` が `nil` でも安全に呼び出せます（protobuf生成コードのgetterは、レシーバがnilの場合ゼロ値を返す設計になっています）。フィールドに直接 `req.Name` でアクセスするより、この形が定石です
- 文字列の連結は `+` でシンプルに書けます: `"Hello, " + req.GetName() + "!"`
