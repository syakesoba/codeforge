//go:build ignore

package main

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type calculatorServer struct {
	UnimplementedCalculatorServer
}

func (s *calculatorServer) Divide(ctx context.Context, req *DivideRequest) (*DivideReply, error) {
	if req.GetDivisor() == 0 {
		return nil, status.Error(codes.InvalidArgument, "divisor must not be zero")
	}
	return &DivideReply{Quotient: req.GetDividend() / req.GetDivisor()}, nil
}

func main() {}
