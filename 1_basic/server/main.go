package main

import (
	"context"
	svc "github.com/kekeee-shine/grpc_training/1_basic/server/service"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/kekeee-shine/grpc_training/1_basic/proto"
	"google.golang.org/grpc"
)

const (
	address = "127.0.0.1"
	port    = ":20051"
)

func main() {

	// listen the tcp port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, svc.NewServer())
	log.Printf("Starting gRPC listener on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func clintMain(wg *sync.WaitGroup) {
	conn, err := grpc.Dial(address+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect :%v", err)
	}

	defer conn.Close()

	c := pb.NewProductInfoClient(conn)
	name := "Apple x"
	description := "just description"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.AddProduct(ctx, &pb.Product{Name: name, Description: description})
	if err != nil {
		log.Fatalf("Could not add product: %v", err)

	}
	log.Printf("Product ID: %s added successfully", r.Value)

	product, err := c.GetProduct(ctx, &pb.ProductID{Value: r.Value})
	if err != nil {
		log.Fatalf("Could not get product: %v", err)

	}
	log.Println("Product", product.String())

	wg.Done()
	log.Printf("Client is closed")
}
