package main

import "google.golang.org/grpc/resolver"

var (
	addrs = []string{"localhost:50051", "localhost:50052"}
)

// Построитель (билдер, builder), который отвечает за создание
// сопоставителя
type exampleResolverBuilder struct{}

func (*exampleResolverBuilder) Build(target resolver.Target,
	cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	// Создание демонстрационного сопоставителя, который находит
	// адреса для lb.example.grpc.io
	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			// lb.example.grpc.io сопоставляется с localhost:50051
			// и localhost:50052
			"lb.example.grpc.io": addrs,
		},
	}
	r.start()
	return r, nil
}

func (*exampleResolverBuilder) Scheme() string { return "exampleScheme" }

// Структура DNS-сопоставителя
type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *exampleResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	_ = r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*exampleResolver) Close()                                  {}

func init() {
	resolver.Register(&exampleResolverBuilder{})
}
