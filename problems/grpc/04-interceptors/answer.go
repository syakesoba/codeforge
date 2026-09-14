//go:build ignore

package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const expectedToken = "Bearer secret-token"

func AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing or invalid token")
	}

	values := md.Get("authorization")
	if len(values) == 0 || values[0] != expectedToken {
		return nil, status.Error(codes.Unauthenticated, "missing or invalid token")
	}

	return handler(ctx, req)
}

func main() {}
