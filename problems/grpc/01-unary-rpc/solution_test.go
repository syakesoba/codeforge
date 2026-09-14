package main

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

// dialServer はメモリ内リスナー（bufconn）でGreeterサーバーを起動し、
// そこに接続したクライアントを返す。実際のTCPポートを使わないため、
// テストが高速かつポート競合の心配なく実行できる。
func dialServer(t *testing.T) GreeterClient {
	t.Helper()

	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	RegisterGreeterServer(s, &server{})
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

func TestSayHello(t *testing.T) {
	client := dialServer(t)

	resp, err := client.SayHello(context.Background(), &HelloRequest{Name: "World"})
	if err != nil {
		t.Fatalf("SayHello returned error: %v", err)
	}
	if want := "Hello, World!"; resp.GetMessage() != want {
		t.Errorf("got %q, want %q", resp.GetMessage(), want)
	}
}

func TestSayHelloDifferentName(t *testing.T) {
	client := dialServer(t)

	resp, err := client.SayHello(context.Background(), &HelloRequest{Name: "Go"})
	if err != nil {
		t.Fatalf("SayHello returned error: %v", err)
	}
	if want := "Hello, Go!"; resp.GetMessage() != want {
		t.Errorf("got %q, want %q", resp.GetMessage(), want)
	}
}
