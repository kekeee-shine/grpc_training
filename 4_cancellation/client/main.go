package main

import (
	"context"
	pb "github.com/kekeee-shine/grpc_training/4_cancellation/proto"
	"google.golang.org/grpc"
	wrapper "google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"log"
	"os"
	"strconv"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	{
		processClient, _ := client.ProcessOrders(ctx)
		for i := 301; i < 302; i++ {

			err := processClient.Send(&wrapper.StringValue{Value: strconv.Itoa(i)})
			if err == io.EOF {
				log.Print("EOF")
				break
			}
		}

		channel := make(chan bool, 1)

		go asyncClientRecvRpc(processClient, channel)

		cancel()

		err := processClient.CloseSend()
		if err != nil {
			return
		}

		<-channel // 阻塞 等待协程中channel写入
	}
}

func asyncClientRecvRpc(processClient pb.OrderManagement_ProcessOrdersClient, ch chan bool) {
	for {
		rlt, err := processClient.Recv()

		if err == io.EOF {
			log.Print("EOF")
			break
		}
		if err != nil {
			log.Printf("Error Receiving messages %v", err)
			break
		} else {
			if rlt != nil {
				log.Println("Process result ", rlt.Value)
			}
		}
	}
	ch <- true
}
