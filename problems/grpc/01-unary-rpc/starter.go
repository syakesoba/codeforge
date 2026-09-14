package main

import "context"

// server は Greeter サービスの実装です。
// UnimplementedGreeterServer を埋め込むことで、将来サービス定義に
// メソッドが追加されても、実装し忘れているだけではコンパイルエラーに
// ならないようにします（gRPCサーバー実装の定石です）。
type server struct {
	UnimplementedGreeterServer
}

// SayHello はクライアントから受け取った名前に挨拶を返す Unary RPC です。
func (s *server) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	// TODO: "Hello, " の後に req.GetName() 、さらに "!" を続けたメッセージを
	// 持つ &HelloReply{Message: ...} を返す
	return nil, nil
}

func main() {}
