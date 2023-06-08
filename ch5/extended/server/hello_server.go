package main

import (
	"context"

	pb "extended/server/ecommerce"
)

type helloServer struct {
}

func (s *helloServer) SayHello(ctx context.Context,
	request *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: "hello world!",
	}, nil
}
