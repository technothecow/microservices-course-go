package grpcServerImpl

import (
	"context"

	grpc "sn/libraries/proto/posts"
	"sn/posts/internal/usecase"
)

func (s *Server) ListPosts(ctx context.Context, req *grpc.ListPostsRequest) (*grpc.ListPostsResponse, error) {
	posts, err := usecase.ListPosts(req)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
