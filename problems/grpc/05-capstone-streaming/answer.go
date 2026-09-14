//go:build ignore

package main

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type counterServer struct {
	UnimplementedCounterServer
}

func (s *counterServer) CountUp(req *CountRequest, stream Counter_CountUpServer) error {
	if req.GetTo() < 1 {
		return status.Error(codes.InvalidArgument, "to must be at least 1")
	}

	for i := int32(1); i <= req.GetTo(); i++ {
		if err := stream.Send(&CountReply{Value: i}); err != nil {
			return err
		}
	}

	return nil
}

func main() {}
