package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	// Импортируем пакет со сгенерированным кодом, который мы только то создали
	// с помощью компилятора protobuf.
	pb "productinfo/service/ecommerce"
)

const (
	port = ":50051"
)

func main() {
	// TCP-слушатель, к которому мы хотим привязать gRPC-сервер,
	// создается на порте 50051.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Для создания нового экземпляра gRPC-сервера применяются API gRPC на Go.
	s := grpc.NewServer()
	// Используя сгенерированный API, регистрируем реализованный ранее сервис
	// на только что созданном gRPC-сервере.
	pb.RegisterProductInfoServer(s, &server{})

	log.Printf("Starting gRPC listener on port " + port)
	// Начинаем прослушивать входящие сообщения на порте 50051.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
