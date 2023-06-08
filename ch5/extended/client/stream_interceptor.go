package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

// Обертка для интерфейса grpc.ClientStream
type wrappedStream struct {
	grpc.ClientStream
}

func newWrappedStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStream{s}
}

// RecvMsg перехватывает сообщения, принимаемые в рамках потокового RPC
func (w *wrappedStream) RecvMsg(m any) error {
	log.Printf("====== [Client Stream Interceptor] "+
		"Receive a message (Type: %T) at %v",
		m, time.Now().Format(time.RFC3339))
	return w.ClientStream.RecvMsg(m)
}

// SendMsg перехватывает сообщения, отправляемые в рамках потокового RPC
func (w *wrappedStream) SendMsg(m any) error {
	log.Printf("====== [Client Stream Interceptor] "+
		"Send a message (Type: %T) at %v",
		m, time.Now().Format(time.RFC3339))
	return w.ClientStream.SendMsg(m)
}

func clientStreamInterceptor(ctx context.Context,
	desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
	streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	// На этапе предобработки есть доступ к RPC-запросу перед
	// его отправкой на сервер
	log.Println("======= [Client Interceptor] ", method)

	// Вызов переданной функции streamer для получения ClientStream
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}

	// Создание обертки вокруг интерфейса ClientStream, перегрузка его
	// методов с использованием логики перехвата и возвращение его
	// клиентскому приложению
	return newWrappedStream(s), nil
}
