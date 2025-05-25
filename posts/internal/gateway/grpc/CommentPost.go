package grpcServerImpl

import (
	"context"
	"log"
	"sn/libraries/postgres"

	grpc "sn/libraries/proto/posts"
	"sn/posts/internal/usecase"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) CommentPost(ctx context.Context, req *grpc.CommentPostRequest) (*emptypb.Empty, error) {
	log.Printf("received create comment request: %v", req.UserId)

	err := usecase.CreateComment(req, postgres.GetPostgresConnection())

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
