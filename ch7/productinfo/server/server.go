package main

import (
	"context"
	"sync"

	"github.com/gofrs/uuid"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"

	pb "productinfo/protobuf"
)

type server struct {
	mu         sync.Mutex
	productMap map[string]*pb.Product
}

func (s *server) AddProduct(ctx context.Context, product *pb.Product) (*wrappers.StringValue, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Добавление новой метрики с помощью созданного ранее счетчика
	customMetricCounter.WithLabelValues(product.Name).Inc()
	out, err := uuid.NewV4()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while generating Product ID: %v", err)
	}
	product.Id = out.String()
	if s.productMap == nil {
		s.productMap = make(map[string]*pb.Product)
	}
	s.productMap[product.Id] = product
	return &wrappers.StringValue{Value: product.Id}, status.New(codes.OK, "").Err()
}

// GetProduct реализует ecommerce.GetProduct
func (s *server) GetProduct(ctx context.Context, productID *wrappers.StringValue) (*pb.Product, error) {
	// Начинаем новый интервал, указывая его имя и контекст
	ctx, span := trace.StartSpan(ctx, "ecommerce.GetProduct")
	// Закончив работу, завершаем интервал
	defer span.End()
	value, exists := s.productMap[productID.Value]
	if exists {
		return value, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Product does not exist. %s", productID.Value)
}
