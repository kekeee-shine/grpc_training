package main

import (
	pb "github.com/kekeee-shine/grpc_training/5_multiplexing/proto"
	svc "github.com/kekeee-shine/grpc_training/5_multiplexing/server/service"
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

	// 在gRPC orderMgtServer上注册订单管理服务
	pb.RegisterOrderManagementServer(s, svc.NewOrderServer())

	// 在gRPC HelloServer上注册问候服务
	pb.RegisterHelloServer(s, svc.NewHelloServer())

	log.Printf("Starting gRPC listener on port " + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
