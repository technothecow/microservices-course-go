package gateway_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	gen "sn/gateway/generated"
	"sn/gateway/internal/usecase"
	stats "sn/libraries/proto/stats"
	mocks "sn/libraries/proto/stats/mocks"
)

var validDate = timestamppb.New(time.Now())

func TestGetPostStats(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStatsClient := mocks.NewMockStatsServiceClient(ctrl)

	mockStatsClient.EXPECT().GetPostStats(gomock.Any(), gomock.Any()).Return(&stats.PostStatsResponse{
		Views:        1,
		Likes:        2,
		Comments:     3,
	}, nil)

	usecase.GetStatsClient = func() (stats.StatsServiceClient, context.Context, func(), error) {
		return mockStatsClient, context.Background(), func() {}, nil
	}

	comments, err := usecase.GetPostStats(&gen.PostStatsRequest{
		PostId: uuid.MustParse(validUUID),
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), comments.Views)
	assert.Equal(t, int64(2), comments.Likes)
	assert.Equal(t, int64(3), comments.Comments)
}

func TestGetPostDynamics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockStatsServiceClient(ctrl)

	mockUsecase.EXPECT().GetPostDynamics(gomock.Any(), gomock.Any()).Return(&stats.PostDynamicsResponse{
		Views: []*stats.PostDynamicsItem{
			{Date: validDate, Value: 1},
			{Date: validDate, Value: 2},
		},
		Likes: []*stats.PostDynamicsItem{
			{Date: validDate, Value: 3},
			{Date: validDate, Value: 4},
		},
		Comments: []*stats.PostDynamicsItem{
			{Date: validDate, Value: 5},
			{Date: validDate, Value: 6},
		},
	}, nil)

	usecase.GetStatsClient = func() (stats.StatsServiceClient, context.Context, func(), error) {
		return mockUsecase, context.Background(), func() {}, nil
	}

	comments, err := usecase.GetPostDynamics(&gen.PostDynamicsRequest{
		PostId: uuid.MustParse(validUUID),
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, comments.Views[0].Value)
	assert.Equal(t, 2, comments.Views[1].Value)
	assert.Equal(t, 3, comments.Likes[0].Value)
	assert.Equal(t, 4, comments.Likes[1].Value)
	assert.Equal(t, 5, comments.Comments[0].Value)
	assert.Equal(t, 6, comments.Comments[1].Value)
}

func TestGetTopPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockStatsServiceClient(ctrl)

	mockUsecase.EXPECT().GetTopPosts(gomock.Any(), gomock.Any()).Return(&stats.TopPostsResponse{
		TopPosts: []*stats.TopPostsItem{
			{PostId: validUUID, Views: 1, Likes: 2, Comments: 3},
			{PostId: validUUID, Views: 4, Likes: 5, Comments: 6},
		},
	}, nil)

	usecase.GetStatsClient = func() (stats.StatsServiceClient, context.Context, func(), error) {
		return mockUsecase, context.Background(), func() {}, nil
	}

	comments, err := usecase.GetTopPosts(&gen.TopPostsRequest{
		Param: gen.TopPostsRequestParamViews,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, comments.Posts[0].Views)
	assert.Equal(t, 2, comments.Posts[0].Likes)
	assert.Equal(t, 3, comments.Posts[0].Comments)
	assert.Equal(t, 4, comments.Posts[1].Views)
	assert.Equal(t, 5, comments.Posts[1].Likes)
	assert.Equal(t, 6, comments.Posts[1].Comments)
}

func TestGetTopUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockStatsServiceClient(ctrl)

	mockUsecase.EXPECT().GetTopUsers(gomock.Any(), gomock.Any()).Return(&stats.TopUsersResponse{
		TopUsers: []*stats.TopUsersItem{
			{UserId: validUUID, Views: 1, Likes: 2, Comments: 3},
			{UserId: validUUID, Views: 4, Likes: 5, Comments: 6},
		},
	}, nil)

	usecase.GetStatsClient = func() (stats.StatsServiceClient, context.Context, func(), error) {
		return mockUsecase, context.Background(), func() {}, nil
	}

	comments, err := usecase.GetTopUsers(&gen.TopUsersRequest{
		Param: gen.TopUsersRequestParamViews,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, comments.Users[0].Views)
	assert.Equal(t, 2, comments.Users[0].Likes)
	assert.Equal(t, 3, comments.Users[0].Comments)
	assert.Equal(t, 4, comments.Users[1].Views)
	assert.Equal(t, 5, comments.Users[1].Likes)
	assert.Equal(t, 6, comments.Users[1].Comments)
}