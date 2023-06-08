package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "productinfo_basic_auth/server/ecommerce"
)

var (
	address  = "localhost:50051"
	hostname = "localhost"
	certFile = "server.crt"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile(certFile, hostname)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	// Инициализируем переменную auth с помощью корректных учетных данных
	// (имени пользователя и пароля). Мы будем использовать хранящиеся
	// в ней значения
	auth := basicAuth{
		username: "admin",
		password: "admin",
	}

	opts := []grpc.DialOption{
		// Переменная auth соответствует интерфейсу, который функция
		// grpc.WithPerRPCCredentials принимает в качестве параметра.
		// Поэтому мы можем передать auth этой функции
		grpc.WithPerRPCCredentials(auth),
		grpc.WithTransportCredentials(creds),
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProductInfoClient(conn)

	res, err := c.GetProduct(context.Background(), &pb.ProductID{
		Value: "102",
	})
	if err != nil {
		log.Fatalf("faield to get product: %v", err)
	}

	log.Println(res.Price)
}
