package main

import (
	"context"
	"fmt"
	"log"

	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "extended/server/ecommerce"
)

type orderMgtServer struct {
}

func (s *orderMgtServer) AddOrder(ctx context.Context, orderReq *pb.Order) (*pb.OrderID, error) {
	// Некорректный запрос. Нужно сгенерировать ошибку и отправить ее клиенту
	if orderReq.Id == "-1" {
		log.Printf("Order ID is invalid! -> Rerceived Order ID %s", orderReq.Id)
		// Создаем новое состояние с кодом ошибки InvalidArgument
		errorStatus := status.New(codes.InvalidArgument, "Invalid information received")
		// Описываем произошедшее с помощью типа ошибки BadRequest_FieldViolation
		// из google.golang.org/genproto/googleapis/rpc/errdetails
		ds, err := errorStatus.WithDetails(&epb.BadRequest_FieldViolation{
			Field: "ID",
			Description: fmt.Sprintf(
				"Order ID received is not valid %s: %s",
				orderReq.Id, orderReq.Destination),
		})
		if err != nil {
			return nil, errorStatus.Err()
		}
		// Возвращаем сгенерированную ошибку
		return nil, ds.Err()
	}

	// Считываем хеш-таблицу с метаданными из входящего контекста
	// удаленного метода
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Println(md)
	}

	// Создаем заголовок
	header := metadata.Pairs("header-key", "val")
	// Отправляем данные в виде заголовка
	_ = grpc.SendHeader(ctx, header)
	// Создаем заключительный блок
	trailer := metadata.Pairs("trailer-key", "val")
	// Отправляем метаданные вместе с заключительным блоком
	_ = grpc.SetTrailer(ctx, trailer)

	return nil, nil
}

func (s *orderMgtServer) ProcessOrders(stream pb.OrderManagement_ProcessOrdersServer) error {
	// Получаем контекст из потока и считываем из него метаданные
	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		log.Println(md)
	}

	// Создаем заголовок
	header := metadata.Pairs("header-key", "val")
	// Отправляем данные в виде заголовка в потоке
	_ = stream.SendHeader(header)
	// Создаем заключительный блок
	trailer := metadata.Pairs("trailer-key", "val")
	// Отправляем метаданные вместе с заключительным блоком в потоке
	stream.SetTrailer(trailer)

	return nil
}
