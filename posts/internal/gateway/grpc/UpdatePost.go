package grpcServerImpl

import (
	"context"

	grpc "sn/libraries/proto/posts"
	"sn/posts/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdatePost(ctx context.Context, req *grpc.UpdatePostRequest) (*grpc.Post, error) {
	post, err := usecase.UpdatePost(req)
	if err != nil {
		if err == usecase.ErrPostNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return post, nil
}
