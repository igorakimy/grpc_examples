package main

import (
	"context"
	"crypto/x509"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	pb "productinfo_google_cloud_auth/server/ecommerce"
)

var (
	address = "localhost:50051"
)

func main() {
	// Вызываем oauth.NewServiceAccountFromFile в целях инициализации
	// credentials.PerRPCCredentials. Для создания учетных данных
	// нужно предоставить файл с действительным токеном
	perRPC, err := oauth.NewServiceAccountFromFile("service-account.json", "scope")
	if err != nil {
		log.Fatalf("faield to create JWT credentials: %v", err)
	}

	pool, _ := x509.SystemCertPool()

	creds := credentials.NewClientTLSFromCert(pool, "")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	opts := []grpc.DialOption{
		// Настраиваем параметры DialOption WithPerRPCCredentials, чтобы
		// аутентификационный токен применялся ко всем удаленным вызовам
		// в рамках одного соединения,  по аналогии с тем, как мы это
		// делали в уже рассмотренных механизмах аутентификации
		grpc.WithPerRPCCredentials(perRPC),
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
