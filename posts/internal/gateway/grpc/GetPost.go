package grpcServerImpl

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"sn/libraries/kafka"
	grpc "sn/libraries/proto/posts"
	"sn/posts/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ViewEvent struct {
	UserId    string `json:"user_id"`
	PostId    string `json:"post_id"`
	Timestamp int64  `json:"timestamp"`
}

var postsViewsTopic = "posts_views"

func (s *Server) GetPost(ctx context.Context, req *grpc.GetPostRequest) (*grpc.Post, error) {
	post, err := usecase.GetPost(req)
	if err != nil {
		log.Printf("failed to get post: %v", err)
		if errors.Is(err, usecase.ErrPostNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	event := ViewEvent{
		UserId:    req.GetRequesterId(),
		PostId:    post.GetId(),
		Timestamp: post.GetCreatedAt().GetSeconds(),
	}
	json_, err := json.Marshal(&event)
	if err != nil {
		log.Printf("json marshal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	} else {
		kafka.SendMessageAsync(postsViewsTopic, []byte(post.GetId()), json_, false)
	}

	return post, nil
}
