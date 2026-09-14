//go:build ignore

package main

import "context"

type server struct {
	UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Hello, " + req.GetName() + "!"}, nil
}

func main() {}
