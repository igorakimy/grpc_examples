package main

import (
	"context"

	pb "productinfo_mtls/server/ecommerce"
)

type server struct{}

func (s *server) AddProduct(ctx context.Context, product *pb.Product) (*pb.ProductID, error) {
	return &pb.ProductID{Value: "105"}, nil
}

func (s *server) GetProduct(ctx context.Context, productID *pb.ProductID) (*pb.Product, error) {
	return &pb.Product{
		Id:    "104",
		Price: 2000.00,
	}, nil
}
