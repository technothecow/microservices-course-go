package grpcServerImpl

import (
	"context"
	"log"

	"sn/libraries/proto/stats"
	"sn/stats/internal/usecase"
)

func (s *Server) GetTopPosts(ctx context.Context, req *stats.GetTopPostsRequest) (*stats.TopPostsResponse, error) {
	comments, err := usecase.GetTopPosts(req)
	if err != nil {
		log.Printf("Failed to get top posts: %v", err)
		return nil, err
	}
	return comments, nil
}