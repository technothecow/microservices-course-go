package grpcServerImpl

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	grpc "sn/libraries/proto/posts"
	"sn/posts/internal/usecase"
)

func (s *Server) DeletePost(ctx context.Context, req *grpc.DeletePostRequest) (*emptypb.Empty, error) {
	err := usecase.DeletePost(req)
	if err != nil {
		if err == usecase.ErrPostNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
