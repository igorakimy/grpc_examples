package main

import (
	"context"

	pb "productinfo_google_cloud_auth/server/ecommerce"
)

type server struct {
	productMap map[string]*pb.Product
}

func (s *server) AddProduct(ctx context.Context, product *pb.Product) (*pb.ProductID, error) {
	panic("implement me")
}

func (s *server) GetProduct(ctx context.Context, productID *pb.ProductID) (*pb.Product, error) {
	if productID.Value == "102" {
		return &pb.Product{
			Id:    "102",
			Price: 1033.00,
		}, nil
	}
	return nil, nil
}
