package main

import (
	"context"
	pb "github.com/kekeee-shine/grpc_training/5_multiplexing/proto"
	"google.golang.org/grpc"
	wrapper "google.golang.org/protobuf/types/known/wrapperspb"
	"log"
	"os"
	"time"
)

const (
	address = "127.0.0.1:20051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect :%v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			os.Exit(0)
		}
	}(conn)

	// 将订单服务客户端绑定至从tcp连接中
	orderClient := pb.NewOrderManagementClient(conn)

	// 将问候服务客户端绑定至从tcp连接中
	helloClient := pb.NewHelloClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	{
		r, err := orderClient.GetOrder(ctx, &wrapper.StringValue{Value: "101"})
		if err != nil {
			log.Fatalf("Could not get order: %v", err)
		}
		log.Printf("Get OrderMgmtServer successfully %v", r)
	}

	{
		r, err := helloClient.SayHello(ctx, &wrapper.StringValue{Value: "Nice to meet you"})
		if err != nil {
			log.Fatalf("Could not get hello from server: %v", err)
		}
		log.Printf("Get HelloServer successfully %v", r)
	}
}
