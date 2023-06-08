package mock_ecommerce

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"google.golang.org/protobuf/proto"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"

	pb "productinfo/server/ecommerce"
)

type rpcMsg struct {
	msg proto.Message
}

func (r *rpcMsg) Matches(msg any) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m, r.msg)
}

func (r *rpcMsg) String() string {
	return fmt.Sprintf("is %s", r.msg)
}

func TestAddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Создаем макет, который ожидает вызова своих удаленных методов
	mockProdInfoClient := NewMockProductInfoClient(ctrl)

	name := "Samsung S10"
	description := "Samsung Galaxy S10 is the latest smart phone, " +
		"launched in February 2019"
	price := float32(700.0)

	req := &pb.Product{
		Name:        name,
		Description: description,
		Price:       price,
	}
	// Программируем макет
	mockProdInfoClient.
		// Ожидаем вызов метода AddProduct
		EXPECT().AddProduct(gomock.Any(), &rpcMsg{msg: req}).
		// Возвращаем фиктивное значение для ID заказа
		Return(&wrappers.StringValue{Value: "Product: " + name}, nil)

	// Вызываем метод теста, который обращается к удаленному
	// методу клиентской заглушки
	testAddProduct(t, mockProdInfoClient)
}

func testAddProduct(t *testing.T, client pb.ProductInfoClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	name := "Samsung S10"
	description := "Samsung Galaxy S10 is the latest smart phone, launched in February 2019"
	price := float32(700.0)

	r, err := client.AddProduct(ctx, &pb.Product{
		Name:        name,
		Description: description,
		Price:       price,
	})
	if err != nil || r.GetValue() != "Product: Samsung S10" {
		t.Errorf("Mocking failed")
	}

	t.Logf("Res: %s", r.GetValue())
}
