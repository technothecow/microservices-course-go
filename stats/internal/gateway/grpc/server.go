package grpcServerImpl

import (
	grpc "sn/libraries/proto/stats"
)

type Server struct {
	grpc.UnimplementedStatsServiceServer
}
