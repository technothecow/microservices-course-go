package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"sn/libraries/kafka"
	pb "sn/libraries/proto/posts"

	"sn/posts/internal/gateway/grpc"
)

func main() {
	defer kafka.CloseProducer()

	listener, err := net.Listen("tcp", ":50003")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPostServiceServer(grpcServer, &grpcServerImpl.Server{})

	log.Println("gRPC server listening on :50002")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
