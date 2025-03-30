package grpcServerImpl

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"sn/libraries/postgres"
	grpc "sn/libraries/proto/users"

	"sn/users/internal/usecase"
)

func (*Server) AuthenticateUser(ctx context.Context, request *grpc.AuthenticateUserRequest) (*grpc.AuthResponse, error) {
	profile, err := usecase.AuthenticateUser(request.GetUsername(), request.GetPassword(), postgres.GetPostgresConnection())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &grpc.AuthResponse{
		Id: profile.Id,
	}, nil
}
