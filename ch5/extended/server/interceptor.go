package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

// orderUnaryServerInterceptor унарный перехватчик на стороне сервера
func orderUnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// Перед вызовом: здесь вы можете перехватить сообщение до того, как будет
	// вызван соответствующий удаленный метод
	log.Println("======= [Server Interceptor] ", info.FullMethod)

	// Вызов RPC-метода с помощью UnaryHandler
	m, err := handler(ctx, req)

	// После вызова: здесь вы можете обработать ответ, сгенерированный
	// в результате удаленного вызова
	log.Printf("Post Proc Message: %s", m)

	// Возвращение RPC-ответа
	return m, err
}
