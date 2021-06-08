package main

import (
	"context"
	"fmt"
	pb "github.com/kekeee-shine/grpc_training/7_resolver/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	wrapper "google.golang.org/protobuf/types/known/wrapperspb"
	"log"
	"os"
	"time"
)

const (
	myScheme      = "example"
	myServiceName = "kekeee.com"

	address = "127.0.0.1:20051"
)

func main() {
	conn, err := grpc.Dial(fmt.Sprintf("%s:///%s", myScheme, myServiceName), grpc.WithInsecure())
	// "example:///kekeee.com"
	if err != nil {
		log.Fatalf("did not connect :%v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			os.Exit(0)
		}
	}(conn)

	client := pb.NewOrderManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// unary request demo::GetOrder
	{
		r, err := client.GetOrder(ctx, &wrapper.StringValue{Value: "101"})
		if err != nil {
			log.Fatalf("Could not ger order: %v", err)

		}
		log.Printf("GerOrder successfully %v", r)
	}

}

type exampleResolverBuilder struct {
	resolver.Builder
}

// Build implement Builder.Build
func (e exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{
		cc:     cc,
		target: target,
		addressMap: map[string][]string{
			myServiceName: {address},
		},
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

// Scheme implement Builder.Scheme
func (e exampleResolverBuilder) Scheme() string {
	return myScheme
}

type exampleResolver struct {
	cc         resolver.ClientConn
	target     resolver.Target
	addressMap map[string][]string
	resolver.Resolver
}

// ResolveNow implement Resolver.ResolveNow
func (e exampleResolver) ResolveNow(options resolver.ResolveNowOptions) {
	addrStrs := e.addressMap[e.target.Endpoint]

	addrs := make([]resolver.Address, len(addrStrs))

	for i, addr := range addrStrs {
		addrs[i] = resolver.Address{
			Addr: addr,
		}
	}

	err := e.cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		return
	}
}

// Close implement Resolver.Close
func (e exampleResolver) Close() {}

func init() {
	// 注册解析器
	resolver.Register(&exampleResolverBuilder{})
}
