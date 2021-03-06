package main

import (
	"context"
	pb "github.com/kekeee-shine/grpc_training/1_basic/proto"
	"google.golang.org/grpc"
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

	log.Printf("Client is closed")
}
