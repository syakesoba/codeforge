package main

// counterServer は Counter サービスの実装です。
type counterServer struct {
	UnimplementedCounterServer
}

// CountUp は 1 から req.GetTo() まで、1つずつストリームで送信する
// サーバーストリーミングRPCです。
func (s *counterServer) CountUp(req *CountRequest, stream Counter_CountUpServer) error {
	// TODO:
	// - req.GetTo() が 1 未満の場合、
	//   status.Error(codes.InvalidArgument, "to must be at least 1") を返す
	// - 1 から req.GetTo() まで1ずつ増やしながら、
	//   stream.Send(&CountReply{Value: i}) を呼び出す
	//   （送信中にエラーが起きたら、そのエラーをそのまま返す）
	// - すべて送り終えたら nil を返す
	return nil
}

func main() {}
