package main

import (
	"log"
	"time"

	"google.golang.org/grpc"
)

// Обертка для потока grpc.ServerStream
type wrappedStream struct {
	grpc.ServerStream
}

// RecvMsg обрабатывает сообщения, принимаемые с помощью потокового RPC
func (w *wrappedStream) RecvMsg(m any) error {
	log.Printf("====== [Server Stream Interceptor Wrapper] "+
		"Receive a message (Type: %T) as %s",
		m, time.Now().Format(time.RFC3339),
	)
	return w.ServerStream.RecvMsg(m)
}

// SendMsg обрабатывает сообщения, оправляемые с помощью потокового RPC
func (w *wrappedStream) SendMsg(m any) error {
	log.Printf("====== [Server Stream Interceptor Wrapper] "+
		"Send a message (Type: %T) at %v",
		m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

// newWrappedStream создает экземпляр обертки
func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

// orderServerStreamInterceptor является реализацией потокового перехватчика
func orderServerStreamInterceptor(srv any,
	ss grpc.ServerStream, info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	// Этап обработки
	log.Println("====== [Server Stream Interceptor] ", info.FullMethod)
	// Вызов метода потокового RPC с помощью обертки
	err := handler(srv, newWrappedStream(ss))
	if err != nil {
		log.Printf("RPC failed with error %v", err)
	}
	return err
}
