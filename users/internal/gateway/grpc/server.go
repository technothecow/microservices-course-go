package grpcServerImpl

import (
	grpc "sn/libraries/proto/users"
)

type Server struct {
	grpc.UnimplementedUserServiceServer
}
