package gateway_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"

	gen "sn/gateway/generated"
	"sn/gateway/internal/usecase"
	posts "sn/libraries/proto/posts"
	"sn/libraries/proto/posts/mocks"
)

var validUUID = "123e4567-e89b-12d3-a456-426614174000"

func TestCreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostsClient := posts_mocks.NewMockPostServiceClient(ctrl)

	usecase.GetPostsClient = func() (posts.PostServiceClient, context.Context, func(), error) {
		return mockPostsClient, context.Background(), func() {}, nil
	}

	t.Run("create post", func(t *testing.T) {
		mockPostsClient.EXPECT().CreatePost(gomock.AssignableToTypeOf(context.Background()), gomock.AssignableToTypeOf(&posts.CreatePostRequest{})).Return(&posts.Post{Id: validUUID}, nil)

		post, err := usecase.CreatePost("userId", &gen.CreatePostRequest{
			Title:       "Test Post",
			Description: "Test Description",
			IsPrivate:   "true",
			Tags:        []string{"test", "post"},
		})

		assert.NoError(t, err)
		assert.Equal(t, validUUID, post.Id.String())
	})

	t.Run("get post", func(t *testing.T) {
		mockPostsClient.EXPECT().GetPost(gomock.Any(), gomock.Any()).Return(&posts.Post{Id: validUUID}, nil)

		post, err := usecase.GetPost(validUUID, "userId")
		assert.NoError(t, err)
		assert.Equal(t, validUUID, post.Id.String())
	})

	t.Run("list posts", func(t *testing.T) {
		mockPostsClient.EXPECT().ListPosts(gomock.Any(), gomock.Any()).Return(&posts.ListPostsResponse{Posts: []*posts.Post{{Id: validUUID}}}, nil)

		posts_, err := usecase.ListPosts("userId", &gen.PaginatedPostsRequest{
			Page:     0,
			Pagesize: 10,
			Tags:     &[]string{"test", "post"},
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(posts_.Posts))
		assert.Equal(t, validUUID, posts_.Posts[0].Id.String())
	})

	t.Run("edit post", func(t *testing.T) {
		mockPostsClient.EXPECT().UpdatePost(gomock.Any(), gomock.Any()).Return(&posts.Post{Id: validUUID}, nil)

		post, err := usecase.EditPost("userId", &gen.EditPostRequest{
			Id:          uuid.MustParse(validUUID),
			Title:       "New Title",
			Description: "New Description",
		})
		assert.NoError(t, err)
		assert.Equal(t, validUUID, post.Id.String())
	})

	t.Run("delete post", func(t *testing.T) {
		mockPostsClient.EXPECT().DeletePost(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, nil)

		err := usecase.DeletePost("userId", validUUID)
		assert.NoError(t, err)
	})

	t.Run("post comment", func(t *testing.T) {
		mockPostsClient.EXPECT().CommentPost(gomock.Any(), gomock.Any()).Return(&emptypb.Empty{}, nil)

		err := usecase.CreateComment("userId", &gen.CreateCommentRequest{
			PostId: uuid.MustParse(validUUID),
			Text:   "Test Comment",
		})
		assert.NoError(t, err)
	})

	t.Run("list post comments", func(t *testing.T) {
		mockPostsClient.EXPECT().ListComments(gomock.Any(), gomock.Any()).Return(&posts.ListCommentsResponse{Comments: []*posts.Comment{{Id: validUUID, PostId: validUUID, UserId: validUUID}}}, nil)

		comments, err := usecase.GetCommentsList(validUUID, &gen.PaginatedCommentsRequest{
			PostId:   uuid.MustParse(validUUID),
			Page:     0,
			Pagesize: 10,
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(comments.Comments))
		assert.Equal(t, validUUID, comments.Comments[0].Id.String())
	})
}
