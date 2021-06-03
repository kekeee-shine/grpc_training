package main

import (
	pb "github.com/kekeee-shine/grpc_training/4_cancellation/proto"
	svc "github.com/kekeee-shine/grpc_training/4_cancellation/server/service"
	"google.golang.org/grpc"
	"log"
	"net"
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
	pb.RegisterOrderManagementServer(s, svc.NewServer())
	log.Printf("Starting gRPC listener on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
