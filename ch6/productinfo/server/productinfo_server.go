package main

import (
	"context"

	pb "productinfo/server/ecommerce"
)

type server struct{}

func (s *server) AddProduct(ctx context.Context, product *pb.Product) (*pb.ProductID, error) {
	panic("implement me")
}

func (s *server) GetProduct(ctx context.Context, productID *pb.ProductID) (*pb.Product, error) {
	if productID.Value == "102" {
		return &pb.Product{
			Id:    "102",
			Price: 3002.00,
		}, nil
	}
	return nil, nil
}
