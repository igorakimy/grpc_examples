package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "productinfo_mtls/server/ecommerce"
)

var (
	port     = ":50051"
	certFile = "server.crt"
	keyFile  = "server.key"
	caFile   = "ca.crt"
)

func main() {
	// Создаем пары ключей X.509 непосредственно из ключа
	// и сертификата сервера
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	// Генерируем пул сертификатов в удостоверяющем центре
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile(caFile)
	if err != nil {
		log.Fatalf("could not read ca certificate: %s", err)
	}

	// Добавляем клиентские сертификаты из удостоверяющего
	// центра в сгенерированный пул
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("failed to append ca certificate")
	}

	// Включаем TLS для всех входящих соединений путем создания
	// аутентификационных данных
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{cert},
			ClientCAs:    certPool,
		})),
	}

	// Создаем новый экземпляр gRPC-сервера, передавая ему
	// аутентификационные данные
	s := grpc.NewServer(opts...)
	// Регистрируем сервис на только что созданном gRPC-сервере,
	// используя сгенерированные API
	pb.RegisterProductInfoServer(s, &server{})
	// Начинаем прослушивать TCP-соединение на порте 50051
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Привязываем gRPC-сервер к слушателю и ждем появления
	// сообщений на порте 50051
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
