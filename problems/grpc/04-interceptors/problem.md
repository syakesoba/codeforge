# Lesson 4: インターセプターで共通処理を挟む

## 今回学ぶこと

- gRPCの「インターセプター」とは何か（Web APIコースのミドルウェアに相当する仕組み）
- リクエストに付随する「メタデータ」の読み取り方
- 認証チェックをインターセプターとして実装する方法

### インターセプター = gRPC版ミドルウェア

「Web APIの基礎」コースや「Gin」コースで、リクエストの前後に共通処理（ロギングや認証）を挟む**ミドルウェア**を学びました。gRPCにも同じ役割を持つ仕組みがあり、**インターセプター**と呼ばれます。

Unary RPC用のインターセプターは、次のような関数シグネチャです。

```go
type UnaryServerInterceptor func(
	ctx context.Context,
	req any,
	info *UnaryServerInfo,
	handler UnaryHandler,
) (resp any, err error)
```

`handler` が「本来呼ばれるはずだったRPCハンドラ」です。インターセプターの中で `handler(ctx, req)` を呼べば処理を継続でき、呼ばなければそこで処理を打ち切れます（ミドルウェアで `next()` を呼ぶかどうかで制御するのと同じ考え方です）。

サーバー作成時に登録すると、**すべてのRPC呼び出し**の前段に自動的に挟まります。

```go
s := grpc.NewServer(grpc.UnaryInterceptor(AuthInterceptor))
```

### メタデータ（Metadata）

gRPCには、リクエスト本体（protobufメッセージ）とは別に、HTTPのヘッダーに近い「メタデータ」という付随情報の仕組みがあります。認証トークンなどは、リクエストメッセージのフィールドではなく、このメタデータに載せて送るのが一般的です。

サーバー側でメタデータを読み取るには `metadata.FromIncomingContext` を使います。

```go
import "google.golang.org/grpc/metadata"

md, ok := metadata.FromIncomingContext(ctx)
if !ok {
	// メタデータが無い
}
values := md.Get("authorization") // []string（同じキーが複数あることもあるため）
```

## 演習

`AuthInterceptor` を実装してください。すでに定義済みの `expectedToken` 定数（`"Bearer secret-token"`）と比較します。

**`AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)`**
1. `metadata.FromIncomingContext(ctx)` でメタデータを取り出す。取り出せなかった場合（`ok` が `false`）は `status.Error(codes.Unauthenticated, "missing or invalid token")` を返す
2. `md.Get("authorization")` の結果が空、または最初の値が `expectedToken` と一致しない場合も同じエラーを返す
3. 一致すれば `handler(ctx, req)` を呼び出し、その戻り値をそのまま返す

## ヒント

- `google.golang.org/grpc/codes`、`google.golang.org/grpc/metadata`、`google.golang.org/grpc/status` のインポートが必要です（「インポートを自動修正」が使えます）
- `md.Get(key)` は `[]string` を返します。存在しないキーなら空スライスが返るので、`len(values) == 0` で判定できます
