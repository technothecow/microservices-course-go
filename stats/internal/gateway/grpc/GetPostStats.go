package grpcServerImpl

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"sn/libraries/proto/stats"
	"sn/stats/internal/usecase"
)

func (s *Server) GetPostStats(ctx context.Context, req *stats.GetPostStatsRequest) (*stats.PostStatsResponse, error) {
	comments, err := usecase.GetPostStats(req)
	if err != nil {
		if errors.Is(err, usecase.ErrPostNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Printf("Failed to get post stats: %v", err)
		return nil, err
	}
	return comments, nil
}