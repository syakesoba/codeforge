package main

import "context"

// calculatorServer は Calculator サービスの実装です。
type calculatorServer struct {
	UnimplementedCalculatorServer
}

// Divide は dividend を divisor で割った商を返す Unary RPC です。
func (s *calculatorServer) Divide(ctx context.Context, req *DivideRequest) (*DivideReply, error) {
	// TODO:
	// - req.GetDivisor() が 0 の場合、codes.InvalidArgument のステータスコードで
	//   エラーメッセージ "divisor must not be zero" を返す
	//   （status.Error(codes.InvalidArgument, "divisor must not be zero") を使う）
	// - それ以外の場合、req.GetDividend() / req.GetDivisor() を計算した商を持つ
	//   &DivideReply{Quotient: ...} と nil を返す
	return nil, nil
}

func main() {}
