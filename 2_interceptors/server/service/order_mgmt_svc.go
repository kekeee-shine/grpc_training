package service

import (
	"context"
	pb "github.com/kekeee-shine/grpc_training/2_interceptors/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	wrapper "google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"log"
	"strings"
)

type Server struct {
	orderMap map[string]*pb.Order
	pb.OrderManagementServer
}

func NewServer() *Server {
	orderMap := make(map[string]*pb.Order)
	orderMap["101"] = &pb.Order{Id: "101", Items: []string{"Google Pixel 3A", "Mac Book Pro"}, Destination: "Mountain View, CA", Price: 1800.00}
	orderMap["102"] = &pb.Order{Id: "102", Items: []string{"Apple Watch S4"}, Destination: "San Jose, CA", Price: 400.00}
	orderMap["103"] = &pb.Order{Id: "103", Items: []string{"Google Home Mini", "Google Nest Hub"}, Destination: "Mountain View, CA", Price: 400.00}
	orderMap["104"] = &pb.Order{Id: "104", Items: []string{"Amazon Echo"}, Destination: "San Jose, CA", Price: 30.00}
	orderMap["105"] = &pb.Order{Id: "105", Items: []string{"Amazon Echo", "Apple iPhone XS"}, Destination: "Mountain View, CA", Price: 30.00}
	return &Server{orderMap: orderMap}
}

//	GetOrder implements proto.OrderManagementServer
func (s Server) GetOrder(ctx context.Context, value *wrapper.StringValue) (*pb.Order, error) {
	log.Println("handle GetOrder request : ", value.GetValue())
	order, exists := s.orderMap[value.Value]
	if exists {
		return order, status.New(codes.OK, "").Err()
	}
	return nil, status.Newf(codes.NotFound, "order %v is not found", value.String()).Err()
}

//	SearchOrders implements proto.OrderManagementServer
func (s Server) SearchOrders(value *wrapper.StringValue, server pb.OrderManagement_SearchOrdersServer) error {
	log.Println("handle SearchOrders request : ", value.GetValue())
	for key, order := range s.orderMap {
		for _, itemStr := range order.Items {
			if strings.Contains(itemStr, value.Value) {
				// Send the matching orders in a stream
				log.Print("Matching Order Found : "+key, " -> Writing Order to the stream ... ")
				err := server.Send(order)
				if err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

//	UpdateOrders implements proto.OrderManagementServer
func (s Server) UpdateOrders(server pb.OrderManagement_UpdateOrdersServer) error {

	ordersStr := "Updated Order IDs : "
	for {
		order, err := server.Recv()
		log.Printf("handle UpdateOrders request %v : ", order)
		if err == io.EOF {
			// Finished reading the order stream.
			return server.SendAndClose(&wrapper.StringValue{Value: "Orders processed " + ordersStr})
		}
		// Update order

		s.orderMap[order.Id] = order

		ordersStr += order.Id + ", "
	}
}

//	ProcessOrders implements proto.OrderManagementServer
func (s Server) ProcessOrders(stream pb.OrderManagement_ProcessOrdersServer) error {
	var tmpOrderMap = make(map[string]*pb.Order)
	for {
		orderId, err := stream.Recv()
		log.Printf("handle UpdateOrders request %s : ", orderId)
		if err == io.EOF {
			// Client has sent all the messages
			// Send remaining shipments

			log.Println("EOF ", orderId)

			for _, tmpOrder := range tmpOrderMap {
				err := stream.Send(&wrapper.StringValue{Value: tmpOrder.Id})
				if err != nil {
					return nil
				}
			}
			return nil
		}
		if err != nil {
			log.Println(err)
			return err
		}
	}
}
