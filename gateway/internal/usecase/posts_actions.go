package usecase

import (
	"context"
	"errors"
	gen "sn/gateway/generated"
	"sn/libraries/proto/posts"
	"strconv"
	"strings"
	"time"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrPostNotFound = errors.New("post not found")
var ErrPostNotAuthorized = errors.New("not authorized")

// Requires defer cl() to be called after
var GetPostsClient = func() (posts.PostServiceClient, context.Context, func(), error) {
	c, err := grpc.NewClient("posts:50003", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, func() {}, err
	}

	client := posts.NewPostServiceClient(c)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return client, ctx, func() { c.Close(); cancel() }, nil
}

func grpcPostToGenPost(post *posts.Post) *gen.Post {
	uuid, err := uuid.Parse(post.GetId())
	if err != nil {
		return nil
	}
	tags := post.GetTags().GetValues()
	if tags == nil {
		tags = []string{}
	}
	return &gen.Post{
		Id:          uuid,
		Title:       post.Title,
		Description: post.Description,
		IsPrivate:   strconv.FormatBool(post.IsPrivate),
		Tags:        tags,
		CreatedAt:   post.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
		UpdatedAt:   post.UpdatedAt.AsTime().Format("2006-01-02 15:04:05"),
		CreatorId:   post.UserId,
	}
}

func CreatePost(userId string, body *gen.CreatePostRequest) (*gen.Post, error) {
	client, ctx, cl, err := GetPostsClient()
	if err != nil {
		return nil, err
	}
	defer cl()

	request := posts.CreatePostRequest{
		Title:         body.Title,
		Description:   body.Description,
		UserId:        userId,
		IsPrivate:     strings.ToLower(body.IsPrivate) == "true",
		Tags:          &posts.Tags{Values: body.Tags},
	}

	response, err := client.CreatePost(ctx, &request)
	if err != nil {
		log.Printf("failed to create post: %v", err)
		return nil, err
	}

	return grpcPostToGenPost(response), nil
}

func GetPost(postId, userId string) (*gen.Post, error) {
	client, ctx, cl, err := GetPostsClient()
	if err != nil {
		return nil, err
	}
	defer cl()

	request := posts.GetPostRequest{
		Id:            postId,
		RequesterId:   userId,
	}

	response, err := client.GetPost(ctx, &request)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrPostNotFound
		}
		return nil, err
	}
	return grpcPostToGenPost(response), nil
}

func ListPosts(userId string, body *gen.PaginatedPostsRequest) (*gen.PostsList, error) {
	client, ctx, cl, err := GetPostsClient()
	if err != nil {
		return nil, err
	}
	defer cl()

	request := posts.ListPostsRequest{
		RequesterId:   userId,
		PageNumber:    int32(body.Page),
		PageSize:      int32(body.Pagesize),
		Tags:          &posts.Tags{Values: *body.Tags},
	}

	response, err := client.ListPosts(ctx, &request)
	if err != nil {
		return nil, err
	}

	posts := make([]gen.Post, len(response.Posts))
	for i, post := range response.Posts {
		posts[i] = *grpcPostToGenPost(post)
	}
	return &gen.PostsList{Posts: posts}, nil
}

func EditPost(userId string, body *gen.EditPostRequest) (*gen.Post, error) {
	client, ctx, cl, err := GetPostsClient()
	if err != nil {
		return nil, err
	}
	defer cl()

	isPrivate := strings.ToLower(body.IsPrivate) == "true"
	var tags *posts.Tags = nil
	if body.Tags != nil {
		tags = &posts.Tags{Values: body.Tags}
	}

	request := posts.UpdatePostRequest{
		Id:          body.Id.String(),
		RequesterId: userId,
		Title:       &body.Title,
		Description: &body.Description,
		IsPrivate:   &isPrivate,
		Tags:        tags,
	}

	response, err := client.UpdatePost(ctx, &request)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrPostNotFound
		} else if status.Code(err) == codes.PermissionDenied {
			return nil, ErrPostNotAuthorized
		}
		return nil, err
	}

	return grpcPostToGenPost(response), nil
}

func DeletePost(userId string, postId string) error {
	client, ctx, cl, err := GetPostsClient()
	if err != nil {
		return err
	}
	defer cl()

	request := posts.DeletePostRequest{
		Id:          postId,
		RequesterId: userId,
	}

	_, err = client.DeletePost(ctx, &request)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return ErrPostNotFound
		} else if status.Code(err) == codes.PermissionDenied {
			return ErrPostNotAuthorized
		}
		return err
	}
	return nil
}
