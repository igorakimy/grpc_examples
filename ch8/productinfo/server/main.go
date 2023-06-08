package main

import (
	"log"
	"net"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	// Импортируем пакет reflection для доступа к API отражения
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc/reflection"

	"productinfo/pb"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		// Создаем на сервере, цепочку унарных перехватчиков. Перехватчики вызываются
		// в порядке их регистрации с помощью gRPC Middleware
		grpc_opentracing.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_recovery.UnaryServerInterceptor(),
	)))

	pb.RegisterProductInfoServer(s, &server{})
	// Регистрируем отражающий сервис на gRPC-сервере
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
