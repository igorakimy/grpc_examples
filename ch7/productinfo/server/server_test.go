package main

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/test/bufconn"

	pb "productinfo/protobuf"
)

const (
	address = "localhost:50051"
	bufSize = 1024 * 1024
)

var listener *bufconn.Listener

func initGRPCServerHTTP2() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &server{})
	reflection.Register(s)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func getBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}

// initGRPCServerBuffConn инициализирует BufConn.
// Пакет bufconn предоставляет net.Conn, реализованный буфером,
// и связанные с ним функции соединения и прослушивания
func initGRPCServerBuffConn() {
	listener = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &server{})
	reflection.Register(s)
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()
}

// TestServer_AddProduct традиционный тест, который запускает сервер
// и клиент для проверки удаленного метода сервиса
func TestServer_AddProduct(t *testing.T) {
	// Запускаем стандартный gRPC-сервер поверх HTTP/2
	initGRPCServerHTTP2()
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Подключаемся к серверному приложению
	c := pb.NewProductInfoClient(conn)

	name := "Samsung S10"
	description := "Samsung Galaxy S10 is the latest smart phone, launched in February 2019"
	price := float32(700.0)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Вызываем удаленный метода AddProduct
	r, err := c.AddProduct(ctx, &pb.Product{
		Name:        name,
		Description: description,
		Price:       price,
	})
	// Проверяем ответ
	if err != nil {
		t.Fatalf("Could not add product: %v", err)
	}

	if r.Value == "" {
		t.Errorf("Invalid Product ID %s", r.Value)
	}
	log.Printf("Res: %s", r.Value)
}

// TestServer_AddProductBufConn тест, который использует BufConn
func TestServer_AddProductBufConn(t *testing.T) {
	initGRPCServerBuffConn()
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(getBufDialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	name := "Samsung S10"
	description := "Samsung Galaxy S10 is the latest smart phone, launched in February 2019"
	price := float32(700.0)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.AddProduct(ctx, &pb.Product{
		Name:        name,
		Description: description,
		Price:       price,
	})
	if err != nil {
		t.Fatalf("Could not add product: %v", err)
	}
	log.Printf(r.Value)
}
