package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Импортируем пакет со сгенерированным кодом, который мы создали
	// с помощью компилятора protobuf.
	pb "productinfo/client/ecommerce"
)

const (
	address = "localhost:50051"
)

func main() {
	// Устанавливаем соединение с сервером с использованием предоставленного адреса
	// ("localhost:50051"). В данном случае соединение между сервером и клиентом
	// будет незащищенным.
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// Закончив работу, закрываем соединение.
	defer conn.Close()

	// Передаем соединение и создаем экземпляр заглушки, который содержит все
	// удаленные методы, доступные на сервере.
	c := pb.NewProductInfoClient(conn)

	name := "Apple iPhone 11"
	description := `Meet Apple iPhone 11. All-new dual-camera system with
					Ultra Wide and Night mode.`
	price := float32(1000.0)

	// Создаем объект Context, который будет передаваться вместе с удаленными
	// вызовами. Он содержит метаданные, такие как идентификатор конечного
	// пользователя, авторизационные токены и крайний срок обработки запроса.
	// Этот объект будет существовать до тех пор, пока запрос не будет обработан.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Вызываем метод addProduct, передавая ему описание товара. В ответ получим
	// ID новой записи, если все пройдет успешно.
	r, err := c.AddProduct(ctx, &pb.Product{
		Name:        name,
		Description: description,
		Price:       price,
	})
	// В противном случае будет возвращена ошибка.
	if err != nil {
		log.Fatalf("Could not add product: %v", err)
	}
	log.Printf("Product ID: %s added successfully", r.Value)

	// Вызываем метод getProduct, передавая ему ID товара. В ответ получим
	// описание товара, если все пройдет успешно.
	product, err := c.GetProduct(ctx, &pb.ProductID{
		Value: r.Value,
	})
	// В противном случае будет возвращена ошибка.
	if err != nil {
		log.Fatalf("Could not get product: %v", err)
	}
	log.Printf("Product: %s", product.String())
}
