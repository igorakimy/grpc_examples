package main

import (
	"context"
	"log"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Импортируем пакет со сгенерированным кодом обратного прокси-сервера
	"productinfo/pb"
)

var (
	// Указываем URL gRPC-сервера. Убедитесь в отм, что ваш сервер работает
	// корректно и доступен по этому пути.
	grpcServerEndpoint = "localhost:50051"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// Цепочка клиентских унарных перехватчиков
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient()),
	}

	// Регистрируем прокси-обработчик для эндпоинта gRPC-сервера. Во время
	// выполнения мультиплексор сопоставляет HTTP-запрос с шаблоном и вызывает
	// подходящий обработчик
	err := pb.RegisterProductInfoHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		log.Fatalf("Fail to register gRPC gateway service endpoint: %v", err)
	}

	// Начинаем принимать запросы на порте 8081
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Could not setup HTTP endpoint: %v", err)
	}
}
