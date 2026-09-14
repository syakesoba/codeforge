package main

import (
	"context"
	"io"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func dialCounter(t *testing.T) CounterClient {
	t.Helper()

	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	RegisterCounterServer(s, &counterServer{})
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

	return NewCounterClient(conn)
}

func TestCountUp(t *testing.T) {
	client := dialCounter(t)

	stream, err := client.CountUp(context.Background(), &CountRequest{To: 5})
	if err != nil {
		t.Fatalf("CountUp returned error: %v", err)
	}

	var got []int32
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Recv returned error: %v", err)
		}
		got = append(got, reply.GetValue())
	}

	want := []int32{1, 2, 3, 4, 5}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i, v := range want {
		if got[i] != v {
			t.Errorf("got[%d] = %d, want %d", i, got[i], v)
		}
	}
}

func TestCountUpSingle(t *testing.T) {
	client := dialCounter(t)

	stream, err := client.CountUp(context.Background(), &CountRequest{To: 1})
	if err != nil {
		t.Fatalf("CountUp returned error: %v", err)
	}

	reply, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv returned error: %v", err)
	}
	if reply.GetValue() != 1 {
		t.Errorf("got %d, want 1", reply.GetValue())
	}

	if _, err := stream.Recv(); err != io.EOF {
		t.Errorf("2件目にio.EOF以外が返りました: %v", err)
	}
}

func TestCountUpInvalid(t *testing.T) {
	client := dialCounter(t)

	stream, err := client.CountUp(context.Background(), &CountRequest{To: 0})
	if err != nil {
		t.Fatalf("CountUp returned error: %v", err)
	}

	_, err = stream.Recv()
	if err == nil {
		t.Fatal("to=0 なのにエラーがnilでした")
	}
	if got := status.Code(err); got != codes.InvalidArgument {
		t.Errorf("status.Code(err) = %v, want %v", got, codes.InvalidArgument)
	}
}
