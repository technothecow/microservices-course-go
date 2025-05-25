package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	gen "sn/gateway/generated"
	"sn/libraries/proto/stats"
)

// Requires defer cl() to be called after
var GetStatsClient = func() (stats.StatsServiceClient, context.Context, func(), error) {
	c, err := grpc.NewClient("stats:50004", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, func() {}, err
	}

	client := stats.NewStatsServiceClient(c)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return client, ctx, func() { c.Close(); cancel() }, nil
}

func GetPostStats(body *gen.PostStatsRequest) (*gen.PostStats, error) {
	client, ctx, cl, err := GetStatsClient()
	if err != nil {
		return nil, err
	}
	defer cl()

	request := stats.GetPostStatsRequest{
		PostId: body.PostId.String(),
	}

	response, err := client.GetPostStats(ctx, &request)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return &gen.PostStats{
		Comments: response.Comments,
		Likes:    response.Likes,
		Views:    response.Views,
	}, nil
}

func GetPostDynamics(body *gen.PostDynamicsRequest) (*gen.PostDynamics, error) {
	client, ctx, cl, err := GetStatsClient()
	if err != nil {
		return nil, err
	}
	defer cl()

	request := stats.GetPostDynamicsRequest{
		PostId: body.PostId.String(),
	}

	response, err := client.GetPostDynamics(ctx, &request)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	genResponse := gen.PostDynamics{}
	genResponse.Views = make([]gen.Dynamic, len(response.GetViews()))
	genResponse.Comments = make([]gen.Dynamic, len(response.GetViews()))
	genResponse.Likes = make([]gen.Dynamic, len(response.GetViews()))

	for i, v := range response.GetViews() {
		genResponse.Views[i] = gen.Dynamic{
			Date:  v.GetDate().AsTime().Format("2006-01-02"),
			Value: int(v.GetValue()),
		}
	}
	for i, v := range response.GetLikes() {
		genResponse.Likes[i] = gen.Dynamic{
			Date:  v.GetDate().AsTime().Format("2006-01-02"),
			Value: int(v.GetValue()),
		}
	}
	for i, v := range response.GetComments() {
		genResponse.Comments[i] = gen.Dynamic{
			Date:  v.GetDate().AsTime().Format("2006-01-02"),
			Value: int(v.GetValue()),
		}
	}

	return &genResponse, nil
}

func GetTopPosts(body *gen.TopPostsRequest) (*gen.TopPosts, error) {
	client, ctx, cl, err := GetStatsClient()
	if err != nil {
		return nil, err
	}
	defer cl()

	var param stats.TopRequestParam
	if body.Param == gen.TopPostsRequestParamViews {
		param = stats.TopRequestParam_VIEWS
	} else if body.Param == gen.TopPostsRequestParamLikes {
		param = stats.TopRequestParam_LIKES
	} else if body.Param == gen.TopPostsRequestParamComments {
		param = stats.TopRequestParam_COMMENTS
	} else {
		return nil, errors.New("invalid param")
	}

	request := stats.GetTopPostsRequest{
		Param: stats.TopRequestParam(param),
	}

	response, err := client.GetTopPosts(ctx, &request)
	if err != nil {
		return nil, err
	}

	genResponse := gen.TopPosts{}
	genResponse.Posts = make([]gen.TopPostsItem, len(response.GetTopPosts()))

	for i, v := range response.GetTopPosts() {
		genResponse.Posts[i] = gen.TopPostsItem{
			PostId: uuid.MustParse(v.GetPostId()),
			Views:  int(v.GetViews()),
			Likes:  int(v.GetLikes()),
			Comments: int(v.GetComments()),
		}
	}

	return &genResponse, nil
}

func GetTopUsers(body *gen.TopUsersRequest) (*gen.TopUsers, error) {
	client, ctx, cl, err := GetStatsClient()
	if err != nil {
		return nil, err
	}
	defer cl()

	var param stats.TopRequestParam
	if body.Param == gen.TopUsersRequestParamViews {
		param = stats.TopRequestParam_VIEWS
	} else if body.Param == gen.TopUsersRequestParamLikes {
		param = stats.TopRequestParam_LIKES
	} else if body.Param == gen.TopUsersRequestParamComments {
		param = stats.TopRequestParam_COMMENTS
	} else {
		return nil, errors.New("invalid param")
	}

	request := stats.GetTopUsersRequest{
		Param: stats.TopRequestParam(param),
	}

	response, err := client.GetTopUsers(ctx, &request)
	if err != nil {
		return nil, err
	}

	genResponse := gen.TopUsers{}
	genResponse.Users = make([]gen.TopUsersItem, len(response.GetTopUsers()))

	for i, v := range response.GetTopUsers() {
		genResponse.Users[i] = gen.TopUsersItem{
			UserId: uuid.MustParse(v.GetUserId()),
			Views:  int(v.GetViews()),
			Likes:  int(v.GetLikes()),
			Comments: int(v.GetComments()),
		}
	}

	return &genResponse, nil
}
