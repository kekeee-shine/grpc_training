package main

import (
	"context"
	pb "github.com/kekeee-shine/grpc_training/2_interceptors/proto"
	"google.golang.org/grpc"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
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

	c := pb.NewOrderManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetOrder(ctx, &wrapperspb.StringValue{Value: "101"})
	if err != nil {
		log.Fatalf("Could not ger order: %v", err)

	}
	log.Printf("GerOrder successfully %v", r)

	//product, err := c.GetProduct(ctx, &pb.ProductID{Value: r.Value})
	//if err != nil {
	//	log.Fatalf("Could not get product: %v", err)
	//
	//}
	//log.Println("Product", product.String())
	//
	//log.Printf("Client is closed")
}
