package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"

	"github.com/golang/protobuf/ptypes/wrappers"

	pb "order_management/server/ecommerce"
)

type server struct {
	orderMap            map[string]*pb.Order
	combinedShipmentMap map[string]*pb.CombinedShipment
}

func (s *server) GetOrder(ctx context.Context, orderId *wrappers.StringValue) (*pb.Order, error) {
	// Реализация сервиса.
	ord := s.orderMap[orderId.Value]
	return ord, nil
}

func (s *server) SearchOrders(
	searchQuery *wrappers.StringValue,
	stream pb.OrderManagement_SearchOrdersServer,
) error {

	for key, order := range s.orderMap {
		log.Print(key, order)
		for _, itemStr := range order.Items {
			log.Print(itemStr)
			// Находим подходящие заказы.
			if strings.Contains(itemStr, searchQuery.Value) {
				// Отправляем подходящий заказ в поток.
				if err := stream.Send(order); err != nil {
					// Проверяем, не возникли ли какие-либо ошибки при потоковой передаче
					// сообщений клиенту.
					return fmt.Errorf("error sending message to stream: %v", err)
				}
				log.Print("Matching Order Found: " + key)
				break
			}
		}
	}
	return nil
}

func (s *server) UpdateOrders(stream pb.OrderManagement_UpdateOrdersServer) error {
	ordersStr := "Updated Order IDs: "
	for {
		// Читаем сообщение из клиентского потока.
		order, err := stream.Recv()
		// Проверяем, не закончился ли поток.
		if err == io.EOF {
			// Завершаем чтение потока заказов.
			return stream.SendAndClose(&wrappers.StringValue{
				Value: "Orders processed " + ordersStr,
			})
		}
		// Обновляем заказ.
		s.orderMap[order.Id] = order

		log.Print("Order ID ", order.Id, ": Updated")
		ordersStr += order.Id + ", "
	}
}

func (s *server) ProcessOrders(stream pb.OrderManagement_ProcessOrdersServer) error {
	for {
		// Читаем ID заказов из входящего потока.
		orderId, err := stream.Recv()
		// Продолжаем читать, пока не обнаружим конец потока.
		if err == io.EOF {
			// При обнаружении конца потока отправляем клиенту все сгруппированные
			// заказы, которые еще остались.
			for _, comb := range s.combinedShipmentMap {
				stream.Send(comb)
			}
			// Сервер завершает поток, возвращая nil.
			return nil
		}
		if err != nil {
			return err
		}
		log.Print("Order ID ", orderId)
		// Логика для объединения заказов в партии
		// на основе адреса доставки.
		batchMarker := rand.Intn(10)
		orderBatchSize := 10
		// Заказы обрабатываются группами. Когда достигается предельный размер
		// партии, все объединенные заказы возвращаются клиенту в виде потока.
		if batchMarker == orderBatchSize {
			// Передаем клиенту поток заказов, объединенных в партии.
			for _, comb := range s.combinedShipmentMap {
				// Передаем клиенту партию объединенных заказов обратно в поток.
				stream.Send(comb)
			}
			batchMarker = 0
			s.combinedShipmentMap = make(map[string]*pb.CombinedShipment)
		} else {
			batchMarker++
		}
	}
}
