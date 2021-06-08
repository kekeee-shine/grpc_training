package main

import (
	"context"
	pb "github.com/kekeee-shine/grpc_training/3_deadlines/proto"
	"google.golang.org/grpc"
	wrapper "google.golang.org/protobuf/types/known/wrapperspb"
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

	// deadline type
	//ctx, cancel := context.WithDeadline(context.Background(),  time.Now().Add(time.Duration(1 * time.Second)))
	// timeout type
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// unary request demo::GetOrder
	{
		r, err := client.GetOrder(ctx, &wrapper.StringValue{Value: "101"})
		if err != nil {
			log.Fatalf("Could not ger order: %v", err)

		}
		log.Printf("GerOrder successfully %v", r)
	}

	//
	//// stream request demo::SearchOrders
	//{
	//	searchStream, _ := client.SearchOrders(ctx, &wrapper.StringValue{Value: "Google"})
	//	for {
	//		searchOrder, err := searchStream.Recv()
	//
	//		if err == io.EOF {
	//			log.Print("EOF")
	//			break
	//		}
	//
	//		if err == nil {
	//			log.Print("Search Result : ", searchOrder)
	//		}
	//	}
	//}
	//
	////	stream request demo::UpdateOrders
	//{
	//	updateClient, _ := client.UpdateOrders(ctx)
	//	for i := 201; i < 205; i++ {
	//		err := updateClient.Send(&pb.Order{Id: strconv.Itoa(i)})
	//		if err != nil {
	//			log.Print(err)
	//		}
	//	}
	//	rlt, _ := updateClient.CloseAndRecv()
	//	log.Println("Update result ", rlt.Value)
	//}
	//
	////	stream request demo::ProcessOrders
	//{
	//	processClient, _ := client.ProcessOrders(ctx)
	//	for i := 301; i < 302; i++ {
	//
	//		err := processClient.Send(&wrapper.StringValue{Value: strconv.Itoa(i)})
	//		if err == io.EOF {
	//			log.Print("EOF")
	//			break
	//		}
	//	}
	//	err := processClient.CloseSend()
	//	if err != nil {
	//		return
	//	}
	//
	//	for {
	//		rlt, err := processClient.Recv()
	//		if err == io.EOF {
	//			log.Print("EOF")
	//			break
	//		}
	//		if rlt != nil {
	//			log.Println("Process result ", rlt.Value)
	//		}
	//	}
	//
	//}
}
