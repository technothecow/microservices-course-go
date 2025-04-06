package grpcServerImpl

import (
	"context"
	"log"
	"sn/libraries/postgres"

	grpc "sn/libraries/proto/posts"
	"sn/posts/internal/usecase"
)

func (s *Server) ListComments(_ context.Context, req *grpc.ListCommentsRequest) (*grpc.ListCommentsResponse, error) {
	log.Printf("received list comments request: %v", req.UserId)

	result, err := usecase.ListComments(req, postgres.GetPostgresConnection())

	if err != nil {
		return nil, err
	}

	return &grpc.ListCommentsResponse{Comments: result}, nil
}
