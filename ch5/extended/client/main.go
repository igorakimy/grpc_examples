package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "extended/client/ecommerce"
)

const (
	address = "localhost:50052"
)

func main() {
	conn, err := grpc.Dial(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),

		// Установление соединения с сервером путем передачи унарного
		// перехватчика в качестве параметра метода Dial
		grpc.WithUnaryInterceptor(orderUnaryClientInterceptor),

		// Регистрация потокового перехватчика
		grpc.WithStreamInterceptor(clientStreamInterceptor),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderManagementClient(conn)

	clientDeadline := time.Now().Add(2 * time.Second)
	// Устанавливаем двухсекундный крайний срок в текущем контексте
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	defer cancel()

	// добавляем заказ
	order1 := pb.Order{
		Id:          "101",
		Items:       []string{"iPhone XS", "Max Book Pro"},
		Destination: "San Jose, CA",
		Price:       2300.00,
	}
	// Вызываем удаленный метод AddOrder и перехватываем любые
	// потенциальные ошибки внутри addErr
	res, addErr := client.AddOrder(ctx, &order1)

	if addErr != nil {
		// Определяем код ошибки с помощью пакета status
		got := status.Code(addErr)
		// Если вызов превысит крайний срок, то возвращается ошибка
		// типа DEADLINE_EXCEEDED
		log.Printf("Error Occured -> addOrder: %v", got)
	} else {
		log.Printf("AddOrder Response -> %v", res.Value)
	}

	// Задействует то же соединение с целью создать клиент сервиса Hello
	helloClient := pb.NewGreeterClient(conn)
	helloResponse, err := helloClient.SayHello(ctx, &pb.HelloRequest{
		Name: "gRPC Up and Running!",
	})
	log.Println(helloResponse)

	// Получаем ссылку на cancel
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

	// Инициируем потоковый RPC
	streamProcOrder, _ := client.ProcessOrders(ctx)
	// Отправляем сервису поток сообщений
	_ = streamProcOrder.Send(&wrappers.StringValue{Value: "102"})
	_ = streamProcOrder.Send(&wrappers.StringValue{Value: "103"})
	_ = streamProcOrder.Send(&wrappers.StringValue{Value: "104"})

	channel := make(chan bool, 1)

	go asyncClientBidirectionalRPC(streamProcOrder, channel)
	time.Sleep(time.Millisecond * 1000)

	// Отменяем (прерываем) удаленный вызов на клиентской стороне
	cancel()
	// Состояние текущего контекста
	log.Printf("RPC Status: %s", ctx.Err())

	_ = streamProcOrder.Send(&wrappers.StringValue{Value: "101"})
	_ = streamProcOrder.CloseSend()

	<-channel

	// Обработка ошибок на стороне клиента

	// Этот заказ недействителен
	order1 = pb.Order{
		Id:          "-1",
		Items:       []string{"iPhone XS", "Mac Book Pro"},
		Destination: "San Jose, CA",
		Price:       2300.00,
	}
	// Вызываем удаленный метод AddOrder и присваиваем ошибку
	// переменной addOrderError
	res, addOrderError := client.AddOrder(ctx, &order1)

	if addOrderError != nil {
		// Получаем код ошибки из пакета grpc/status
		errorCode := status.Code(addOrderError)
		// Сравниваем код ошибки с InvalidArgument
		if errorCode == codes.InvalidArgument {
			log.Printf("Invalid Argument Error: %s", errorCode)
			// Получаем состояние
			errorStatus := status.Convert(addOrderError)
			for _, d := range errorStatus.Details() {
				switch info := d.(type) {
				// Проверяем, имеет ли ошибка тип BadRequest_FieldViolation
				case *epb.BadRequest_FieldViolation:
					log.Printf("Request Field Invalied: %s", info)
				default:
					log.Printf("Unexpected error type: %s", info)
				}
			}
		} else {
			log.Printf("Unhandled error: %s", errorCode)
		}
	} else {
		log.Print("AddOrder Response -> ", res.Value)
	}

	// Создаем метаданные
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"kn", "vn",
	)
	// Создаем новый контекст с новыми метаданными
	mdCtx := metadata.NewOutgoingContext(context.Background(), md)

	// Добавляем дополнительные метаданные в существующий контекст
	ctxA := metadata.AppendToOutgoingContext(mdCtx,
		"k1", "v1", "k1", "v2", "k2", "v3")

	// Унарный RPC на основе нового контекста с метаданными
	_, err = client.AddOrder(ctxA, &order1)

	// Тот же контекст можно использовать и для потокового RPC
	_, err = client.ProcessOrders(ctxA)

	// Чтение метаданных gRPC на стороне клиента

	// Переменная для хранения заголовка и заключительного блока,
	// возвращенных удаленным вызовом
	var header, trailer metadata.MD

	// Передаем ссылку на заголовок и заключительный блок, чтобы
	// сохранить значения для унарного RPC
	_, err = client.AddOrder(
		ctx,
		&order1,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)

	stream, err := client.ProcessOrders(ctx)

	// Получаем заголовки из потока
	header, err = stream.Header()

	// Получаем из потока заключительные блоки, которые используются
	// для отправки кодов и сообщений о состоянии
	trailer = stream.Trailer()

	// Балансировка нагрузки на стороне клиента
	exampleScheme := "example"
	exampleServiceName := "lb.example.grpc.io"
	// Устанавливаем gRPC-соединение с помощью схемы example и имени сервиса.
	// Протокол возвращается сопоставителем в рамках клиентского приложения
	pickfirstConn, err := grpc.Dial(fmt.Sprintf("%s:///%s",
		exampleScheme, exampleServiceName),
		// Указываем алгоритм балансировки нагрузки, который выбирает первый
		// сервер в списке эндпоинтов
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "pick_first": {} } ]}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer pickfirstConn.Close()

	log.Println("==== Calling helloworld.Greeter/SayHello with pick_first ====")
	makeRPCs(pickfirstConn, 10)

	// Устанавливаем еще одно соединение ClientConn с помощью алгоритма round_robin
	roundrobinConn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", exampleScheme, exampleServiceName),
		// Балансируем нагрузки с помощью алгоритма циклического перебора
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [ { "round_robin": {} } ]}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer roundrobinConn.Close()

	log.Println("==== Calling helloworld.Greeter/SayHello with round_robin ====")
	makeRPCs(roundrobinConn, 10)
}

func asyncClientBidirectionalRPC(streamProcOrder pb.OrderManagement_ProcessOrdersClient, c chan bool) {
	combinedShipment, errProcOrder := streamProcOrder.Recv()
	if errProcOrder != nil {
		// Возвращаем ошибку при попытке получить сообщения из отмененного контекста
		log.Printf("Error Receiving message %v", errProcOrder)
	}
	log.Println(combinedShipment)
	<-c
}

func makeRPCs(conn *grpc.ClientConn, amount int) {
	// Логика создания RPC
}
