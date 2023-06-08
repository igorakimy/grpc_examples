package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "extended/server/ecommerce"
)

const (
	port = "50051"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(
		// Регистрация унарного перехватчика на gRPC-сервере
		grpc.UnaryInterceptor(orderUnaryServerInterceptor),

		// Регистрация потокового перехватчика на gRPC-сервере
		grpc.StreamInterceptor(orderServerStreamInterceptor),
	)

	// Регистрируем сервис OrderManagement на gRPC-сервере
	pb.RegisterOrderManagementServer(s, &orderMgtServer{})

	// Регистрируем сервис Hello на gRPC-сервере
	pb.RegisterGreeterServer(s, &helloServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
