package grpcServerImpl

import (
	"context"
	"log"

	grpc "sn/libraries/proto/posts"
	"sn/posts/internal/usecase"
)

func (s *Server) CreatePost(ctx context.Context, req *grpc.CreatePostRequest) (*grpc.Post, error) {
	log.Printf("received create post request: %v", req.UserId)
	post, err := usecase.CreatePost(req)
	if err != nil {
		return nil, err
	}
	return post, nil
}
