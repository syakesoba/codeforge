# Lesson 3: ステータスコードでエラーを返す

## 今回学ぶこと

- gRPC専用のエラー表現「ステータスコード」とは何か
- `status.Error` でステータスコード付きのエラーを作る方法
- HTTPのステータスコードとの違い

### gRPCのステータスコード

これまでのコースの「Web APIの基礎」では、エラーを `http.StatusBadRequest`（400）のようなHTTPステータスコードで表現しました。gRPCにも同じ役割を持つ、専用の**ステータスコード**の体系があります（`google.golang.org/grpc/codes` パッケージ）。

代表的なものをいくつか挙げます。

| コード | 意味 | HTTPでの近い概念 |
|---|---|---|
| `codes.OK` | 成功 | 200 |
| `codes.InvalidArgument` | リクエストの内容が不正 | 400 |
| `codes.NotFound` | 対象が見つからない | 404 |
| `codes.Unauthenticated` | 認証が必要 | 401 |
| `codes.PermissionDenied` | 権限がない | 403 |
| `codes.Internal` | サーバー内部のエラー | 500 |

HTTPのステータスコードより種類は少なく絞られていますが、「クライアントが悪い（4xx相当）」か「サーバーが悪い（5xx相当）」かを機械的に判定しやすいよう設計されています。

### `status.Error` でエラーを作る

サーバー側のハンドラがエラーを返すときは、`google.golang.org/grpc/status` パッケージの `status.Error` を使い、コードとメッセージをセットで返します。

```go
import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

if req.GetDivisor() == 0 {
	return nil, status.Error(codes.InvalidArgument, "divisor must not be zero")
}
```

普通の `error` として扱いつつ、Lesson 2で見た `status.Code(err)` でクライアント側からコードを取り出せます。素の `errors.New("...")` を返すこともできますが、その場合はgRPCが自動的に `codes.Unknown` として扱ってしまい、クライアントが状況に応じた分岐をしづらくなります。何が起きたかをクライアントに正しく伝えるために、意図に合ったコードを明示的に選ぶことが大切です。

## 演習

`Calculator` サービスの実装である `calculatorServer` の `Divide` メソッドを実装してください。

**`(s *calculatorServer) Divide(ctx context.Context, req *DivideRequest) (*DivideReply, error)`**
- `req.GetDivisor()` が `0` の場合、`status.Error(codes.InvalidArgument, "divisor must not be zero")` を返す
- それ以外の場合、`req.GetDividend() / req.GetDivisor()` を計算し、その商を持つ `&DivideReply{Quotient: ...}` と `nil` を返す

## ヒント

- `google.golang.org/grpc/codes` と `google.golang.org/grpc/status` のインポートが必要です（「インポートを自動修正」が使えます）
- Goの整数の割り算はゼロ除算でパニックするため、`Divisor == 0` のチェックを**先に**行う必要があります
