package main

import (
	"context"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// Импортируем пакет со сгенерированным кодом, который мы только
	// что создали с помощью компилятора protobuf.
	pb "productinfo/service/ecommerce"
)

// Структура server является абстракцией сервера и позволяет подключать
// к нему методы сервиса.
type server struct {
	productMap map[string]*pb.Product
}

// AddProduct реализует ecommerce.AddProduct.
//
// Метод AddProduct принимает в качестве параметра Product и возвращает ProductID.
// Структуры Product и ProductID определены в файле product_info.pb.go, который
// был автоматически сгенерирован из определения product_info.proto.
func (s *server) AddProduct(ctx context.Context, product *pb.Product) (*pb.ProductID, error) {
	out, err := uuid.NewV4()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while generating Product ID", err)
	}
	product.Id = out.String()
	if s.productMap == nil {
		s.productMap = make(map[string]*pb.Product)
	}
	s.productMap[product.Id] = product
	return &pb.ProductID{Value: product.Id}, status.New(codes.OK, "").Err()
}

// GetProduct реализует ecommerce.GetProduct.
//
// Метод GetProduct принимает в качестве параметра ProductID и возвращает Product.
//
// У обоих методов также есть параметр Context. Объект Context существует на
// протяжении жизненного цикла запроса и содержит такие метаданные, как
// идентификатор конечного пользователя, авторизационные токены и крайний
// срок обработки запроса.
//
// Помимо итогового значения, оба метода могут возвращать ошибки (методы могут иметь
// несколько возвращаемых типов). Эти ошибки передаются потребителями и могут быть
// обработаны на клиентской стороне.
func (s *server) GetProduct(ctx context.Context, productID *pb.ProductID) (*pb.Product, error) {
	value, exists := s.productMap[productID.Value]
	if exists {
		return value, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Product does not exist.", productID.Value)
}
