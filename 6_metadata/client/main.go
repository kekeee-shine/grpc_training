package main

import (
	"context"
	pb "github.com/kekeee-shine/grpc_training/6_metadata/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

	client := pb.NewOrderManagementClient(conn)
	// 实例metadata
	md := metadata.Pairs("cst_key1", "cst_value1",
		"client_time", time.Now().Format(time.Stamp),
		"server_time", "")
	// 创建含有metadata的上下文 本质是一个context.WithValue（）
	mdCtx := metadata.NewOutgoingContext(context.Background(), md)

	// append
	metadata.AppendToOutgoingContext(mdCtx, "cst_key2", "cst_value2", "cst_key3", "cst_value3")

	// 如需设置超时时间等
	//ctx,_:= context.WithTimeout(mdCtx, time.Second*10)
	var header, trailer metadata.MD

	// unary request demo::GetOrder
	{
		r, err := client.GetOrder(mdCtx, &wrapper.StringValue{Value: "101"}, grpc.Header(&header), grpc.Trailer(&trailer))
		if err != nil {
			log.Fatalf("Could not ger order: %v", err)

		}
		log.Printf("GerOrder successfully %v", r)

		if cTimeMap, ok := header["client_time"]; ok {
			log.Printf("kv from header['client_time']:\n")
			for k, v := range cTimeMap {
				log.Printf("%d. %s\n", k, v)
			}
		} else {
			log.Printf("No client_time in header")
		}

		// 接收header返回中server_time元数据
		if sTimeMap, ok := header["server_time"]; ok {
			log.Printf("kv from header['server_time']:\n")
			for k, v := range sTimeMap {
				log.Printf("%d. %s\n", k, v)
			}
		} else {
			log.Printf("No server_time in header")
		}

		for _, trailerMap := range trailer {
			for i, value := range trailerMap {
				log.Printf("%d. %s\n", i, value)
			}
		}
	}

}
