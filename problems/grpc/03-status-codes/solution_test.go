package main

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func dialCalculator(t *testing.T) CalculatorClient {
	t.Helper()

	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	RegisterCalculatorServer(s, &calculatorServer{})
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

	return NewCalculatorClient(conn)
}

func TestDivideSuccess(t *testing.T) {
	client := dialCalculator(t)

	resp, err := client.Divide(context.Background(), &DivideRequest{Dividend: 10, Divisor: 2})
	if err != nil {
		t.Fatalf("Divide returned error: %v", err)
	}
	if want := int64(5); resp.GetQuotient() != want {
		t.Errorf("got %d, want %d", resp.GetQuotient(), want)
	}
}

func TestDivideNegative(t *testing.T) {
	client := dialCalculator(t)

	resp, err := client.Divide(context.Background(), &DivideRequest{Dividend: -9, Divisor: 3})
	if err != nil {
		t.Fatalf("Divide returned error: %v", err)
	}
	if want := int64(-3); resp.GetQuotient() != want {
		t.Errorf("got %d, want %d", resp.GetQuotient(), want)
	}
}

func TestDivideByZero(t *testing.T) {
	client := dialCalculator(t)

	_, err := client.Divide(context.Background(), &DivideRequest{Dividend: 10, Divisor: 0})
	if err == nil {
		t.Fatal("0で割ったのにエラーがnilでした")
	}
	if got := status.Code(err); got != codes.InvalidArgument {
		t.Errorf("status.Code(err) = %v, want %v", got, codes.InvalidArgument)
	}
}
