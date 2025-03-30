package grpcServerImpl

import (
	grpc "sn/libraries/proto/posts"
)

type Server struct {
	grpc.UnimplementedPostServiceServer
}
