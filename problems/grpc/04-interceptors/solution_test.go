package main

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// echoGreeterServer は AuthInterceptor をテストするための固定実装。
type echoGreeterServer struct {
	UnimplementedGreeterServer
}

func (s *echoGreeterServer) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Hello, " + req.GetName() + "!"}, nil
}

func dialWithInterceptor(t *testing.T) GreeterClient {
	t.Helper()

	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer(grpc.UnaryInterceptor(AuthInterceptor))
	RegisterGreeterServer(s, &echoGreeterServer{})
	go func() { _ = s.Serve(lis) }()
	t.Cleanup(s.Stop)

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	return NewGreeterClient(conn)
}

func TestAuthInterceptorMissingToken(t *testing.T) {
	client := dialWithInterceptor(t)

	_, err := client.SayHello(context.Background(), &HelloRequest{Name: "World"})
	if err == nil {
		t.Fatal("トークンなしで呼び出したのにエラーがnilでした")
	}
	if got := status.Code(err); got != codes.Unauthenticated {
		t.Errorf("status.Code(err) = %v, want %v", got, codes.Unauthenticated)
	}
}

func TestAuthInterceptorWrongToken(t *testing.T) {
	client := dialWithInterceptor(t)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer wrong-token")
	_, err := client.SayHello(ctx, &HelloRequest{Name: "World"})
	if err == nil {
		t.Fatal("誤ったトークンで呼び出したのにエラーがnilでした")
	}
	if got := status.Code(err); got != codes.Unauthenticated {
		t.Errorf("status.Code(err) = %v, want %v", got, codes.Unauthenticated)
	}
}

func TestAuthInterceptorCorrectToken(t *testing.T) {
	client := dialWithInterceptor(t)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", expectedToken)
	resp, err := client.SayHello(ctx, &HelloRequest{Name: "World"})
	if err != nil {
		t.Fatalf("正しいトークンなのにエラーが返りました: %v", err)
	}
	if want := "Hello, World!"; resp.GetMessage() != want {
		t.Errorf("got %q, want %q", resp.GetMessage(), want)
	}
}
