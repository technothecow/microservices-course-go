package gateway_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"github.com/google/uuid"

	gen "sn/gateway/generated"
	"sn/gateway/internal/usecase"
	posts "sn/libraries/proto/posts"
)

var validUUID = "123e4567-e89b-12d3-a456-426614174000"

type MockPostsClient struct {
	mock.Mock
}

func (m *MockPostsClient) CreatePost(ctx context.Context, in *posts.CreatePostRequest, opts ...grpc.CallOption) (*posts.Post, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*posts.Post), args.Error(1)
}

func (m *MockPostsClient) GetPost(ctx context.Context, in *posts.GetPostRequest, opts ...grpc.CallOption) (*posts.Post, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*posts.Post), args.Error(1)
}

func (m *MockPostsClient) ListPosts(ctx context.Context, in *posts.ListPostsRequest, opts ...grpc.CallOption) (*posts.ListPostsResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*posts.ListPostsResponse), args.Error(1)
}

func (m *MockPostsClient) DeletePost(ctx context.Context, in *posts.DeletePostRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func (m *MockPostsClient) UpdatePost(ctx context.Context, in *posts.UpdatePostRequest, opts ...grpc.CallOption) (*posts.Post, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*posts.Post), args.Error(1)
}

func TestCreatePost(t *testing.T) {
	mockPostsClient := new(MockPostsClient)
	mockPostsClient.On("CreatePost", mock.Anything, mock.AnythingOfType("*posts.CreatePostRequest"), mock.AnythingOfType("[]grpc.CallOption")).Return(&posts.Post{Id: validUUID}, nil)
	usecase.GetPostsClient = func() (posts.PostServiceClient, context.Context, func(), error) {
		return mockPostsClient, context.Background(), func() {}, nil
	}

	post, err := usecase.CreatePost("userId", &gen.CreatePostRequest{
		Title:       "Test Post",
		Description: "Test Description",
		IsPrivate:   "true",
		Tags:        []string{"test", "post"},
	})

	assert.NoError(t, err)
	assert.Equal(t, validUUID, post.Id.String())
}

func TestGetPost(t *testing.T) {
	mockPostsClient := new(MockPostsClient)
	mockPostsClient.On("GetPost", mock.Anything, mock.AnythingOfType("*posts.GetPostRequest"), mock.AnythingOfType("[]grpc.CallOption")).Return(&posts.Post{Id: validUUID}, nil)
	usecase.GetPostsClient = func() (posts.PostServiceClient, context.Context, func(), error) {
		return mockPostsClient, context.Background(), func() {}, nil
	}

	post, err := usecase.GetPost(validUUID, "userId")
	assert.NoError(t, err)
	assert.Equal(t, validUUID, post.Id.String())
}

func TestListPosts(t *testing.T) {
	mockPostsClient := new(MockPostsClient)
	mockPostsClient.On("ListPosts", mock.Anything, mock.AnythingOfType("*posts.ListPostsRequest"), mock.AnythingOfType("[]grpc.CallOption")).Return(&posts.ListPostsResponse{Posts: []*posts.Post{{Id: validUUID}}}, nil)
	usecase.GetPostsClient = func() (posts.PostServiceClient, context.Context, func(), error) {
		return mockPostsClient, context.Background(), func() {}, nil
	}

	posts, err := usecase.ListPosts("userId", &gen.PaginatedPostsRequest{
		Page:     0,
		Pagesize: 10,
		Tags:     &[]string{"test", "post"},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(posts.Posts))
	assert.Equal(t, validUUID, posts.Posts[0].Id.String())
}

func TestEditPost(t *testing.T) {
	mockPostsClient := new(MockPostsClient)
	mockPostsClient.On("UpdatePost", mock.Anything, mock.AnythingOfType("*posts.UpdatePostRequest"), mock.AnythingOfType("[]grpc.CallOption")).Return(&posts.Post{Id: validUUID}, nil)
	usecase.GetPostsClient = func() (posts.PostServiceClient, context.Context, func(), error) {
		return mockPostsClient, context.Background(), func() {}, nil
	}

	post, err := usecase.EditPost("userId", &gen.EditPostRequest{
		Id:          uuid.MustParse(validUUID),
		Title:       "New Title",
		Description: "New Description",
	})
	assert.NoError(t, err)
	assert.Equal(t, validUUID, post.Id.String())
}

func TestDeletePost(t *testing.T) {
	mockPostsClient := new(MockPostsClient)
	mockPostsClient.On("DeletePost", mock.Anything, mock.AnythingOfType("*posts.DeletePostRequest"), mock.AnythingOfType("[]grpc.CallOption")).Return(&emptypb.Empty{}, nil)
	usecase.GetPostsClient = func() (posts.PostServiceClient, context.Context, func(), error) {
		return mockPostsClient, context.Background(), func() {}, nil
	}

	err := usecase.DeletePost("userId", validUUID)
	assert.NoError(t, err)
}
