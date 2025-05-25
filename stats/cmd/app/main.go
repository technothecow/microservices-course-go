package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "sn/libraries/proto/stats"

	"sn/stats/internal/gateway/grpc"
	"sn/stats/internal/workers"
)

func main() {
	workerManager := workers.NewManager()
	go workerManager.Start()

	listener, err := net.Listen("tcp", ":50004")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStatsServiceServer(grpcServer, &grpcServerImpl.Server{})

	log.Println("gRPC server listening on :50004")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
