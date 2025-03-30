package grpcServerImpl

import (
	"context"
	"log"

	grpc "sn/libraries/proto/posts"
	"sn/posts/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetPost(ctx context.Context, req *grpc.GetPostRequest) (*grpc.Post, error) {
	post, err := usecase.GetPost(req)
	if err != nil {
		log.Printf("failed to get post: %v", err)
		if err == usecase.ErrPostNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())	
	}
	return post, nil
}
