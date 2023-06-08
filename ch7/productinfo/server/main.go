package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	// Указываем внешние библиотеки, необходимые для включения мониторинга.
	// В экосистеме gRPC уже есть ряд готовых перехватчиков для поддержки
	// Prometheus, и здесь мы будем использовать их
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	// Указываем внешние библиотеки, необходимые для включения мониторинга.
	// Пакет gRPC OpenCensus предоставляем готовые обработчики для поддержки
	// OpenCensus, и здесь мы будем использовать их
	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
	"google.golang.org/grpc"

	pb "productinfo/protobuf"
)

const (
	port = ":50051"
)

var (
	// Создаем реестр метрик. Он станет хранить все данные, которые сборщики
	// регистрируют в системе. Если нужно добавить новый сборщик, то его
	// следует указать в реестре
	reg = prometheus.NewRegistry()

	// Создаем стандартные клиентские метрики, которые определены в библиотеке
	grpcMetrics = grpc_prometheus.NewServerMetrics()

	// Создаем собственные счетчик метрик с именем product_mgt_server_handle_count
	customMetricCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "product_mgt_server_handle_count",
		Help: "Total number of RPCs handled on the server.",
	}, []string{"name"})
)

func init() {
	// Регистрируем стандартные серверные метрики и наш сборщик в реестре
	reg.MustRegister(grpcMetrics, customMetricCounter)
}

func main() {
	// Запускаем сервер z-Pages. Для визуализации метрик служит HTTP-путь
	// на порте 8081 начинающийся с /debug
	go func() {
		mux := http.NewServeMux()
		zpages.Handle(mux, "/debug")
		log.Fatal(http.ListenAndServe("127.0.0.1:8081", mux))
	}()

	// Регистрируем средства экспорта собираемых данных. Добавленный
	// нами объект PrintExporter экспортирует данные в консоль. Это
	// сделано в целях демонстрации; обычно мы не рекомендуем вести
	// журнал для всех промышленных заданий
	view.RegisterExporter(&exporter.PrintExporter{})

	// Регистрируем представления для сбора количества запросов к серверу.
	// Нам доступно несколько готовых представлений, которые собирают байты,
	// полученные/отправленные в ходе каждого удаленного вызова, латентность
	// этих вызовов и информацию о том, сколько из них завершилось. Мы также
	// можем создавать собственные представления для сбора данных
	if err := view.Register(ocgrpc.DefaultServerViews...); err != nil {
		log.Fatal(err)
	}

	// Инициализация трассировки
	initTracing()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем HTTP-сервер для Prometheus. HTTP-путь для сбора метрик
	// начинается с /metrics и находится на порте 9092
	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%d", 9092),
	}

	// Создаем gRPC-сервер с обработчиком статистики и обработчиком метрик.
	// Поскольку наш сервис унарный, мы используем grpcMetrics.UnaryServerInterceptor.
	// Для потоковых сервисов предусмотрен другой перехватчик,
	// grpcMetrics.StreamServerInterceptor
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
	)

	// Регистрируем на gRPC-сервере сервис ProductInfo
	pb.RegisterProductInfoServer(grpcServer, &server{})

	// Инициализируем все стандартные метрики
	grpcMetrics.InitializeMetrics(grpcServer)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	log.Printf("Starting gRPC listener on port " + port)

	// Начинаем прослушивать входящие сообщения на порте 50051
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
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
