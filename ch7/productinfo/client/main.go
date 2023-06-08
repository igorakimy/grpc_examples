package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	// Указываем внешние библиотеки, необходимые для включения мониторинга
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opencensus.io/trace"

	// Указываем внешние библиотеки, необходимые для включения мониторинга
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"

	pb "productinfo/protobuf"
)

const (
	// address = "productinfo:50051"
	address = "localhost:50051"
)

func main() {
	// Регистрируем средства экспорта собранных данных и трассировок.
	// Здесь мы добавляем объект PrintExporter, который записывает
	// экспортированные данные в консоль. Это сделано в целях демонстрации;
	// обычно мы не рекомендуем вести журнал для всех промышленных заданий
	view.RegisterExporter(&exporter.PrintExporter{})

	// Регистрируем представления для сбора количества запросов к серверу.
	// Нам доступно несколько готовых представлений, которые собирают
	// байты, полученные/отправленные в ходе каждого удаленного вызова,
	// латентность этих вызовов и информацию о том, сколько из них
	// завершилось. Мы также можем создавать собственные представления
	// для сбора данных
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		log.Fatal(err)
	}

	// Создаем реестр метрик. Как и в случае с серверным кодом, реестр
	// будет хранить все данные, которые сборщики регистрируют в системе.
	// Если нужно добавить новый сборщик, то его следует указать в реестре
	reg := prometheus.NewRegistry()
	// Создаем стандартные клиентские метрики, которые определены в библиотеке
	grpcMetrics := grpc_prometheus.NewClientMetrics()
	// Регистрируем стандартные серверные метрики и наш сборщик в реестре
	reg.MustRegister(grpcMetrics)

	// Вызываем функцию initTracing, инициализируем средство экспорта
	// Jaeger и регистрируем его в трассировщике
	initTracing()

	// Устанавливаем соединение с сервером
	conn, err := grpc.Dial(address,
		// Добавляем клиентский обработчик статистики
		grpc.WithStatsHandler(&ocgrpc.ClientHandler{}),

		// Добавляем перехватчик метрик. Это унарный клиент, вследствие чего
		// мы используем grpcMetrics.UnaryClientInterceptor. Для потоковых
		// клиентов предусмотрен другой перехватчик, grpcMetrics.StreamClientInterceptor
		grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),

		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// Закончив работу, закрываем соединение
	defer conn.Close()

	// Создаем HTTP-сервер для Prometheus. HTTP-путь для сбора метрик
	// начинается с /metrics и находится на порте 9094
	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%d", 9094),
	}

	// Запускаем HTTP-сервер для Prometheus
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server")
		}
	}()

	// Создаем на основе установленного соединения клиентскую заглушку
	client := pb.NewProductInfoClient(conn)

	// Начинаем новый интервал, указывая его имя и контекст
	ctx, span := trace.StartSpan(context.Background(),
		"ecommerce.ProductInfoClient")

	name := "Apple iPhone 11"
	description := `Meet Apple iPhone 11. All-new dual-camera system with Ultra Wide and Night mode.`
	price := float32(1000.0)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// Вызываем удаленный метод AddProduct, передавая ему описание
	// нового товара
	r, err := client.AddProduct(ctx, &pb.Product{
		Name:        name,
		Description: description,
		Price:       price,
	})

	if err != nil {
		log.Fatalf("Could not add product: %v", err)
	}
	log.Printf("Product ID: %s added successfully", r.Value)

	// Вызываем удаленный метод GetProduct, передавая ему wrappers.StringValue
	product, err := client.GetProduct(ctx, &wrappers.StringValue{
		Value: r.Value,
	})

	if err != nil {
		log.Fatalf("Could not get product: %v", err)
	}
	log.Printf("Product: %s", product.String())
	// Закончив работу, завершаем интервал
	span.End()
}

func initTracing() {
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})
	agentEndpointURI := "localhost:6831"
	collectorEndpointURI := "http://localhost:14268/api/traces"
	// Создаем средство экспорта Jaeger, указывая путь для сбора данных,
	// имя сервиса и эндпоинт
	exp, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: collectorEndpointURI,
		AgentEndpoint:     agentEndpointURI,
		Process: jaeger.Process{
			ServiceName: "product_info",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	// Регистрируем средство экспорта в трассировщике OpenCensus
	trace.RegisterExporter(exp)
}
