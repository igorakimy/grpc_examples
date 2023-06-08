package main

import (
	"context"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"

	"productinfo/pb"
)

type server struct {
	productMap map[string]*pb.Product
}

func (s *server) AddProduct(ctx context.Context, product *pb.Product) (*wrappers.StringValue, error) {
	out, err := uuid.NewV4()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate Product ID")
	}
	product.Id = out.String()
	if s.productMap == nil {
		s.productMap = make(map[string]*pb.Product, 0)
	}
	s.productMap[product.Id] = product
	return &wrappers.StringValue{Value: product.GetId()}, nil
}

func (s *server) GetProduct(ctx context.Context, productID *wrappers.StringValue) (*pb.Product, error) {
	product, exist := s.productMap[productID.Value]
	if exist {
		return product, nil
	}
	return &pb.Product{}, status.Error(codes.InvalidArgument, "Product ID does not exist")
}
