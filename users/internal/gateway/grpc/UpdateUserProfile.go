package grpcServerImpl

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"sn/libraries/postgres"
	grpc "sn/libraries/proto/users"
	"sn/users/internal/usecase"
)

func (*Server) UpdateUserProfile(ctx context.Context, request *grpc.UpdateUserProfileRequest) (*grpc.UserProfileResponse, error) {
	result, err := usecase.UpdateUserProfile(
		request.GetId(),
		request.GetFullName(),
		request.GetDateOfBirth(),
		request.GetPhoneNumber(),
		postgres.GetPostgresConnection())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &grpc.UserProfileResponse{
		Id:        result.Id,
		Username:  result.Username,
		Email:     result.Email,
		CreatedAt: timestamppb.New(result.CreatedAt),
		UpdatedAt: timestamppb.New(result.UpdatedAt),
		LastLogin: timestamppb.New(result.LastLogin),
		IsActive:  result.IsActive,
	}
	if result.FullName != nil {
		response.FullName = *result.FullName
	}
	if result.PhoneNumber != nil {
		response.PhoneNumber = *result.PhoneNumber
	}
	if result.DateOfBirth != nil {
		response.DateOfBirth = *result.DateOfBirth
	}

	return response, nil
}
