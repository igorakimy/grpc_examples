package main

import (
	"crypto/tls"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "productinfo/server/ecommerce"
)

var (
	port     = "50051"
	certFile = "server.crt"
	keyFile  = "server.key"
)

func main() {
	// Считываем и анализируем открытый/закрытый ключи
	// и создаем сертификат, чтобы включить TLS
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load key pair: %v", err)
	}

	// Включаем TLS для всех входящих соединений, используя
	// сертификаты для аутентификации
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}

	// Создаем новый экземпляр gRPC-сервера, передавая
	// ему аутентифицированные данные
	s := grpc.NewServer(opts...)

	// Регистрируем реализованный сервис на только что
	// созданном gRPC-сервере с помощью сгенерированных API
	pb.RegisterProductInfoServer(s, &server{})

	// Начинаем прослушивать TCP-соединение на порте 50051
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Привязываем gRPC-сервер к слушателю и ждем появления
	// сообщений на порте 50051
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
