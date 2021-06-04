package service

import (
	"context"
	pb "github.com/kekeee-shine/grpc_training/5_multiplexing/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	wrapper "google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"log"
	"strings"
)

type OrderServer struct {
	orderMap map[string]*pb.Order
	pb.OrderManagementServer
}

func NewOrderServer() *OrderServer {
	orderMap := make(map[string]*pb.Order)
	orderMap["101"] = &pb.Order{Id: "101", Items: []string{"Google Pixel 3A", "Mac Book Pro"}, Destination: "Mountain View, CA", Price: 1800.00}
	orderMap["102"] = &pb.Order{Id: "102", Items: []string{"Apple Watch S4"}, Destination: "San Jose, CA", Price: 400.00}
	orderMap["103"] = &pb.Order{Id: "103", Items: []string{"Google Home Mini", "Google Nest Hub"}, Destination: "Mountain View, CA", Price: 400.00}
	orderMap["104"] = &pb.Order{Id: "104", Items: []string{"Amazon Echo"}, Destination: "San Jose, CA", Price: 30.00}
	orderMap["105"] = &pb.Order{Id: "105", Items: []string{"Amazon Echo", "Apple iPhone XS"}, Destination: "Mountain View, CA", Price: 30.00}
	return &OrderServer{orderMap: orderMap}
}

//	GetOrder implements proto.OrderManagementServer
func (s OrderServer) GetOrder(ctx context.Context, value *wrapper.StringValue) (*pb.Order, error) {
	log.Println("Handle GetOrder request : ", value.GetValue())
	order, exists := s.orderMap[value.Value]
	if exists {
		return order, status.New(codes.OK, "").Err()
	}
	return nil, status.Newf(codes.NotFound, "order %v is not found", value.String()).Err()
}

//	SearchOrders implements proto.OrderManagementServer
func (s OrderServer) SearchOrders(value *wrapper.StringValue, server pb.OrderManagement_SearchOrdersServer) error {
	log.Println("Handle SearchOrders request : ", value.GetValue())
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
func (s OrderServer) UpdateOrders(server pb.OrderManagement_UpdateOrdersServer) error {

	ordersStr := "Updated Order IDs : "
	for {
		order, err := server.Recv()
		log.Printf("Handle UpdateOrders request %v : ", order)
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
func (s OrderServer) ProcessOrders(stream pb.OrderManagement_ProcessOrdersServer) error {
	tmpIdArray := make([]string, 0)
	for {

		// You can determine whether the current RPC is cancelled by the other party.
		if stream.Context().Err() == context.Canceled {
			log.Printf(" Context Cacelled for this stream: -> %s", stream.Context().Err())
			log.Printf("Stopped processing any more order of this stream!")
			return stream.Context().Err()
		}

		orderId, err := stream.Recv()

		log.Printf("Handle ProcessOrders request %s : ", orderId)
		if err == io.EOF {
			// Client has sent all the messages
			// Send remaining shipments

			log.Println("EOF ", orderId)

			for Id := range tmpIdArray {
				err := stream.Send(&wrapper.StringValue{Value: tmpIdArray[Id]})
				if err != nil {
					return nil
				}
			}
			log.Println("ProcessOrders send end ")
			return nil
		} else {
			tmpIdArray = append(tmpIdArray, orderId.Value)
		}
		if err != nil {
			log.Println(err)
			return err
		}
	}
}
