package main

import (
	"context"
	"net"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// fakeGreeterServer は CallSayHello をテストするための固定実装。
// 空文字の名前を受け取るとエラーを返す。
type fakeGreeterServer struct {
	UnimplementedGreeterServer
}

func (s *fakeGreeterServer) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	return &HelloReply{Message: "Hello, " + req.GetName() + "!"}, nil
}

func dialClient(t *testing.T) GreeterClient {
	t.Helper()

	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	RegisterGreeterServer(s, &fakeGreeterServer{})
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

func TestCallSayHelloSuccess(t *testing.T) {
	client := dialClient(t)

	got, err := CallSayHello(context.Background(), client, "World")
	if err != nil {
		t.Fatalf("CallSayHello returned error: %v", err)
	}
	if want := "Hello, World!"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCallSayHelloError(t *testing.T) {
	client := dialClient(t)

	_, err := CallSayHello(context.Background(), client, "")
	if err == nil {
		t.Fatal("空の名前を渡したのにエラーがnilでした")
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("status.Code(err) = %v, want %v", status.Code(err), codes.InvalidArgument)
	}
	if !strings.Contains(err.Error(), "call SayHello") {
		t.Errorf("エラーメッセージがラップされていません。実際: %q", err.Error())
	}
}
