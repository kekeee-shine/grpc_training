package main

import (
	"context"
	pb "github.com/kekeee-shine/grpc_training/2_interceptors/proto"
	"google.golang.org/grpc"
	wrapper "google.golang.org/protobuf/types/known/wrapperspb"
	"io"
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
	clientDeadline := time.Now().Add(time.Duration(200 * time.Second))
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// unary request demo
	//r, err := client.GetOrder(ctx, &wrapper.StringValue{Value: "101"})
	//if err != nil {
	//	log.Fatalf("Could not ger order: %v", err)
	//
	//}
	//log.Printf("GerOrder successfully %v", r)

	// stream request demo
	searchStream, _ := client.SearchOrders(ctx, &wrapper.StringValue{Value: "Google"})
	for {
		searchOrder, err := searchStream.Recv()

		if err == io.EOF {
			log.Print("EOF")
			break
		}

		if err == nil {
			log.Print("Search Result : ", searchOrder)
		}
	}

	//product, err := c.GetProduct(ctx, &pb.ProductID{Value: r.Value})
	//if err != nil {
	//	log.Fatalf("Could not get product: %v", err)
	//
	//}
	//log.Println("Product", product.String())
	//
	//log.Printf("Client is closed")
}
