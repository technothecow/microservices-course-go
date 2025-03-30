package grpcServerImpl

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"sn/libraries/postgres"
	grpc "sn/libraries/proto/users"

	"sn/users/internal/usecase"
)

func (*Server) DeleteUserProfile(ctx context.Context, req *grpc.DeleteUserProfileRequest) (*emptypb.Empty, error) {
	err := usecase.DeleteUser(req.GetId(), postgres.GetPostgresConnection())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
