package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "productinfo_basic_auth/server/ecommerce"
)

var (
	port               = ":50051"
	certFile           = "server.crt"
	keyFile            = "server.key"
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid credentials")
)

func main() {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load key pair: %v", err)
	}
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),

		// Добавляем новый серверный параметр (grpc.ServerOption) вместе
		// с TLS-сертификатом сервера. Добавляем с помощью вызова
		// grpc.UnaryInterceptor перехватчик, который будет направлять
		// все клиентские запросы к функции ensureValidBasicCredentials
		grpc.UnaryInterceptor(ensureValidBasicCredentials),
	}

	s := grpc.NewServer(opts...)
	pb.RegisterProductInfoServer(s, &server{})

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Basic ")
	return token == base64.StdEncoding.EncodeToString([]byte("admin:admin"))
}

// ensureValidBasicCredentials проверяет подлинность вызывающей стороны.
// Объект context.Context содержит нужные нам метаданные, которые будут
// существовать на протяжении обработки запроса
func ensureValidBasicCredentials(ctx context.Context, req any,
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	// Извлекаем метаданные из контекста, получаем значение поля
	// authorization и проверяем аутентификационные данные. Ключи
	// внутри metadata.MD переводятся в нижний регистр, поэтому
	// при проверке их значений нужно использовать прописные буквы
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	return handler(ctx, req)
}
