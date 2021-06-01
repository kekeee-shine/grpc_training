package service

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/kekeee-shine/grpc_training/1_basic/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type Server struct {
	productMap map[string]*pb.Product
	// UnimplementedProductInfoServer has implemented all mtd of service.ProductInfoServer
	pb.ProductInfoServer
}

func NewServer() *Server {
	return &Server{productMap: nil}
}

//	AddProduct implements service.ProductInfoServer
func (s *Server) AddProduct(ctx context.Context, in *pb.Product) (*pb.ProductID, error) {
	log.Println("Product", in.String())
	out, err := uuid.NewUUID()
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"Error while generating Product ID", err)
	}
	in.Id = out.String()
	if s.productMap == nil {
		s.productMap = make(map[string]*pb.Product)
	}

	s.productMap[in.Id] = in
	return &pb.ProductID{Value: in.Id}, status.New(codes.OK, "").Err()
}

//	GetProduct implements service.ProductInfoServer
func (s *Server) GetProduct(ctx context.Context, in *pb.ProductID) (*pb.Product, error) {
	value, exist := s.productMap[in.Value]
	if exist {
		return value, status.New(codes.OK, "").Err()
	}
	return value, status.New(codes.NotFound, "").Err()
}
