package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func orderUnaryClientInterceptor(
	ctx context.Context, method string, req, reply any,
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// На этапе предобработки есть доступ к RPC-запросу
	// перед его отправкой на сервер
	log.Println("Method: " + method)

	// Вызов RPC-метода с помощью UnaryInvoker
	err := invoker(ctx, method, req, reply, cc, opts...)

	// Этап постобработки, где можно обработать ответ
	// или возникшую ошибку
	log.Println(reply)

	// Возвращение ошибки клиентскому gRPC-приложению вместе
	// с ответом, который передается в виде аргумента
	return err
}
