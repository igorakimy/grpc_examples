package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/igorakimy/example1"
)

type server struct {
}

// AddProduct удаленный метод для добавления товара
func (s *server) AddProduct(ctx context.Context, in *pb.Product) (*pb.ProductID, error) {
	// бизнес-логика
	return nil, nil
}

// GetProduct удаленный для получения товара
func (s *server) GetProduct(ctx context.Context, in *pb.ProductID) (*pb.Product, error) {
	// бизнес-логика
	return nil, nil
}

func main() {
	// Запуск gRPC-сервера для сервиса ProductInfo

	lis, _ := net.Listen("tcp", ":8000")
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
