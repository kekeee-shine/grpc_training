package main

import (
	pb "github.com/kekeee-shine/grpc_training/3_deadlines/proto"
	svc "github.com/kekeee-shine/grpc_training/3_deadlines/server/service"
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

	//s := grpc.NewServer(grpc.UnaryInterceptor(interceptors.OrderUnaryServerInterceptor1),
	//	grpc.ChainUnaryInterceptor(interceptors.OrderUnaryServerInterceptor2, interceptors.OrderUnaryServerInterceptor3))
	s := grpc.NewServer()
	pb.RegisterOrderManagementServer(s, svc.NewServer())
	log.Printf("Starting gRPC listener on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
