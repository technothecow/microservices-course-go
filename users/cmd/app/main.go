package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "sn/libraries/proto/users"

	"sn/users/internal/gateway/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50002")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &grpcServerImpl.Server{})

	log.Println("gRPC server listening on :50002")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
