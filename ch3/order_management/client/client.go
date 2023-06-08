package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "order_management/client/ecommerce"
)

const (
	address = ":50051"
)

func main() {
	// Устанавливаем соединение с сервером.
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewOrderManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Получаем заказ.
	retrievedOrder, err := c.GetOrder(ctx, &wrappers.StringValue{
		Value: "106",
	})
	if err != nil {
		log.Fatalf("failed to retrieve order: %v", err)
	}
	log.Print("GetOrder Response -> : ", retrievedOrder)

	// Функция SearchOrders возвращает клиентский поток OrderManagement_SearchOrdersClient,
	// у которого есть метод Recv().
	searchStream, _ := c.SearchOrders(ctx, &wrappers.StringValue{
		Value: "Google",
	})

	for {
		// Вызываем метод Recv() из клиентского потока для последовательного
		// получения ответов типа Order.
		searchOrder, err := searchStream.Recv()
		// При обнаружении конца потока Recv возвращает io.EOF.
		if err == io.EOF {
			break
		}
		// Обрабатываем другие потенциальные ошибки.
		log.Print("Search Result: ", searchOrder)
	}

	// Вызов удаленного метода UpdateOrders.
	updateStream, err := c.UpdateOrders(ctx)
	// Обработка ошибок, связанных с UpdateOrders.
	if err != nil {
		log.Fatalf("%v.UpdateOrders(_) = _, %v", c, err)
	}

	updOrder1 := pb.Order{
		Items:       []string{"apples", "banana"},
		Description: "Desc 1",
		Price:       100,
		Destination: "France",
	}

	updOrder2 := pb.Order{
		Items:       []string{"apples", "banana"},
		Description: "Desc 1",
		Price:       100,
		Destination: "Spain",
	}

	updOrder3 := pb.Order{
		Items:       []string{"apples", "banana"},
		Description: "Desc 1",
		Price:       100,
		Destination: "Germany",
	}

	// Отправка обновленных заказов в клиентский поток.

	// Обновляем заказ 1.
	if err := updateStream.Send(&updOrder1); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateStream, &updOrder1, err)
	}

	// Обновляем заказ 2.
	if err := updateStream.Send(&updOrder2); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateStream, &updOrder2, err)
	}

	// Обновляем заказ 3.
	if err := updateStream.Send(&updOrder3); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateStream, &updOrder3, err)
	}

	// Закрытие потока и получение ответа.
	updateRes, err := updateStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", updateStream, err, nil)
	}
	log.Printf("Update Orders Res: %s", updateRes)

	// Обработка заказов.
	// Вызываем удаленный метод и получаем ссылку на поток для записи и чтения на клиентской стороне.
	streamProcOrder, _ := c.ProcessOrders(ctx)

	// Отправляем сообщение сервису.
	if err := streamProcOrder.Send(&wrappers.StringValue{Value: "102"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", c, "102", err)
	}

	if err := streamProcOrder.Send(&wrappers.StringValue{Value: "103"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", c, "103", err)
	}

	if err := streamProcOrder.Send(&wrappers.StringValue{Value: "104"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", c, "104", err)
	}

	// Создаем канал для горутин.
	channel := make(chan struct{})
	// Вызываем функцию с помощью горутин, чтобы распараллелить чтение
	// сообщений, возвращаемых сервисом.
	go asyncClientBidirectionalRPC(streamProcOrder, channel)
	// Имитируем задержку при отправке сервису некоторых сообщений.
	time.Sleep(time.Millisecond * 1000)

	if err := streamProcOrder.Send(&wrappers.StringValue{Value: "101"}); err != nil {
		log.Fatalf("%v.Send(%v) = %v", c, "101", err)
	}

	// Сигнализируем о завершении клиентского потока (с ID заказов).
	if err := streamProcOrder.CloseSend(); err != nil {
		log.Fatal(err)
	}

	<-channel
}

func asyncClientBidirectionalRPC(streamProcOrder pb.OrderManagement_ProcessOrdersClient, c chan struct{}) {
	for {
		// Читаем сообщения сервиса на клиентской стороне.
		combinedShipment, errProcOrder := streamProcOrder.CloseAndRecv()
		// Условие для обнаружения конца потока.
		if errProcOrder == io.EOF {
			break
		}
		log.Printf("Combined shipment: %v", combinedShipment.OrderList)
	}
	<-c
}
