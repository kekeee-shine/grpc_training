package service

import (
	"context"
	pb "github.com/kekeee-shine/grpc_training/5_multiplexing/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	wrapper "google.golang.org/protobuf/types/known/wrapperspb"
)

type HelloServer struct {
	pb.HelloServer
}

func NewHelloServer() *HelloServer {
	return &HelloServer{}
}

func (h HelloServer) SayHello(ctx context.Context, value *wrapperspb.StringValue) (*wrapperspb.StringValue, error) {
	return &wrapper.StringValue{Value: "Nice to meet you,too!"}, status.Newf(codes.OK, "").Err()
}
