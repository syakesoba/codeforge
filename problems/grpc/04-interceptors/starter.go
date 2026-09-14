package main

import (
	"context"

	"google.golang.org/grpc"
)

// expectedToken は、このレッスンで使う固定の認証トークンです
// （実運用では環境変数やシークレット管理サービスから読み込みます）。
const expectedToken = "Bearer secret-token"

// AuthInterceptor は、リクエストのメタデータに正しいトークンが
// 含まれている場合だけ次のハンドラを呼び出す UnaryServerInterceptor です。
func AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// TODO:
	// 1. metadata.FromIncomingContext(ctx) でメタデータを取り出す
	//    （取り出せなかった場合は status.Error(codes.Unauthenticated, "missing or invalid token") を返す）
	// 2. "authorization" キーの値（md.Get("authorization")）が1件以上あり、
	//    その最初の値が expectedToken と一致するか確認する。
	//    一致しなければ status.Error(codes.Unauthenticated, "missing or invalid token") を返す
	// 3. 一致すれば handler(ctx, req) を呼び出し、その結果をそのまま返す
	return nil, nil
}

func main() {}
