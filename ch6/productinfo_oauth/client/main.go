package main

import (
	"context"
	"log"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"

	pb "productinfo_oauth/server/ecommerce"
)

var (
	address  = "localhost:50051"
	hostname = "localhost"
	certFile = "server.crt"
)

func main() {
	// OAuth требует, чтобы транспортный канал был защищен
	creds, err := credentials.NewClientTLSFromFile(certFile, hostname)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	// Подготавливаем учетные данные для соединения. Но прежде
	// необходимо предоставить значение токена OAuth2. Здесь мы
	// используем строку, прописанную в коде
	auth := oauth.NewOauthAccess(fetchToken())

	opts := []grpc.DialOption{
		// Указываем один и тот же токен OAuth в параметрах всех
		// вызовов в рамках одного соединения. Если вы хотите
		// указывать токен для каждого вызова отдельно, то
		// используйте CallOption
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

func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "some-secret-token",
	}
}
