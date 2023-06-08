package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "productinfo_mtls/server/ecommerce"
)

var (
	address  = "localhost:50051"
	hostname = "localhost"
	certFile = "client.crt"
	keyFile  = "client.key"
	caFile   = "ca.crt"
)

func main() {
	// Создаем пары ключей X.509 непосредственно из ключа
	// и сертификата сервера
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("could not load client key pair: %s", err)
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
		log.Fatalf("failed to append ca certificates")
	}

	// Указываем транспортные аутентификационные данные в виде
	// параметров соединения. Поле ServerName должно быть равно
	// значению Common Name, указанному в сертификате
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			ServerName:   hostname,
			Certificates: []tls.Certificate{cert},
			RootCAs:      certPool,
		})),
	}

	// Устанавливаем безопасное соединение с сервером
	// и передаем параметры
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// Закончив работу, закрываем соединение
	defer conn.Close()

	// Передаем соединение и создаем экземпляр заглушки,
	// который содержит все удаленные методы, доступные
	// на сервере
	c := pb.NewProductInfoClient(conn)

	res, err := c.GetProduct(context.Background(), &pb.ProductID{Value: "102"})
	if err != nil {
		log.Fatalf("product not found: %v", err)
	}

	log.Println(fmt.Sprintf("Product price: %f", res.Price))
}
