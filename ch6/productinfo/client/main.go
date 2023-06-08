package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "productinfo/client/ecommerce"
)

var (
	address  = "localhost:50051"
	hostname = "localhost"
	certFile = "server.crt"
)

func main() {
	// Считываем и анализируем публичный сертификат, чтобы включить TLS
	creds, err := credentials.NewClientTLSFromFile(certFile, hostname)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	// Указываем аутентификационные данные для транспортного протокола
	// с помощью DialOption
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	// Устанавливаем безопасное соединение с сервером, передавая
	// параметры аутентификации
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// Закрываем соединение после окончания работы
	defer conn.Close()

	// Передаем соединение с создаем заглушку. Ее экземпляр содержит
	// все удаленные методы, которые можно вызывать на сервере
	c := pb.NewProductInfoClient(conn)

	res, err := c.GetProduct(context.Background(), &pb.ProductID{
		Value: "102",
	})
	if err != nil {
		log.Fatalf("faield to get product: %v", err)
	}

	log.Println(res.Price)
}
