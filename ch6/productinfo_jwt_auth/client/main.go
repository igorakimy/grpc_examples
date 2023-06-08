package main

import (
	"context"
	"log"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	pb "productinfo_jwt_auth/server/ecommerce"
)

var (
	address  = "localhost:50051"
	hostname = "localhost"
	certFile = "server.crt"
)

func main() {
	// Вызываем oauth.NewJWTAccessFromFile с целью инициализации
	// credentials.PerRPCCredentials. Для создания учетных данных
	// нужно предоставить файл с действительным токеном
	jwtCreds, err := oauth.NewJWTAccessFromFile("token.json")
	if err != nil {
		log.Fatalf("failed to create JWT credentials: %v", err)
	}

	creds, err := credentials.NewClientTLSFromFile(certFile, hostname)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	opts := []grpc.DialOption{
		// Настраиваем параметры DialOption WithPerRPCCredentials,
		// чтобы токен JWT применялся ко всем удаленным вызовам в
		// рамках одного соединения
		grpc.WithPerRPCCredentials(jwtCreds),
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

func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "some-secret-token",
	}
}
