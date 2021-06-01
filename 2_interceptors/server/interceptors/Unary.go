package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"strings"
)

// OrderUnaryServerInterceptor1 Server :: Unary Interceptor
// log interceptor
func OrderUnaryServerInterceptor1(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Pre-processing logic
	// Gets info about the current RPC call by examining the args passed in
	log.Println("======= [Server Interceptor 1] ", info.FullMethod)
	log.Printf(" Pre Proc Message : %s", req)

	// Invoking the handler to complete the normal execution of a unary RPC.
	m, err := handler(ctx, req)

	// Post processing logic
	log.Printf(" Post Proc Message : %s", m)
	log.Println("======= [Server Interceptor 1] End ")
	return m, err
}

// OrderUnaryServerInterceptor2 Server :: Unary Interceptor
// id check interceptor
func OrderUnaryServerInterceptor2(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Pre-processing logic
	// Gets info about the current RPC call by examining the args passed in
	log.Println("-------[Server Interceptor 2] ", info.FullMethod)
	if !strings.Contains(info.FullMethod, "get") {
		log.Println("just accept the get request")
	}
	log.Println("this is a get request, pass")
	// Invoking the handler to complete the normal execution of a unary RPC.
	m, err := handler(ctx, req)
	log.Println("------- [Server Interceptor 1] End ")
	return m, err
}

// OrderUnaryServerInterceptor3 Server::Unary Interceptor
// do other things interceptor
func OrderUnaryServerInterceptor3(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("******* [Server Interceptor 3] ", info.FullMethod)
	// Invoking the handler to complete the normal execution of a unary RPC.
	m, err := handler(ctx, req)
	log.Println("******* [Server Interceptor 3] End ")
	return m, err
}
