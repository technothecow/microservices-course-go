package grpcServerImpl

import (
	"context"
	"encoding/json"
	"log"

	"sn/libraries/kafka"
	grpc "sn/libraries/proto/posts"
	"sn/posts/internal/usecase"
)

func (s *Server) CreatePost(ctx context.Context, req *grpc.CreatePostRequest) (*grpc.Post, error) {
	log.Printf("received create post request: %v", req.UserId)
	post, err := usecase.CreatePost(req)
	if err != nil {
		return nil, err
	}

	event := ViewEvent{
		UserId:    req.GetUserId(),
		PostId:    post.GetId(),
		Timestamp: post.GetCreatedAt().GetSeconds(),
	}
	json_, err := json.Marshal(&event)
	if err != nil {
		log.Printf("json marshal error: %v", err)
	} else {
		kafka.SendMessageAsync(postsViewsTopic, []byte(post.GetId()), json_, false)
	}

	return post, nil
}
