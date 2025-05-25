package grpcServerImpl

import (
	"context"
	"log"

	"sn/libraries/proto/stats"
	"sn/stats/internal/usecase"
)

func (s *Server) GetTopUsers(ctx context.Context, req *stats.GetTopUsersRequest) (*stats.TopUsersResponse, error) {
	comments, err := usecase.GetTopUsers(req)
	if err != nil {
		log.Printf("Failed to get top users: %v", err)
		return nil, err
	}
	return comments, nil
}