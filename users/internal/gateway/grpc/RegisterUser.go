package grpcServerImpl

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/timestamppb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"sn/libraries/kafka"
	"sn/libraries/postgres"

	"sn/users/internal/usecase"

	grpc "sn/libraries/proto/users"
)

var usersCreatedTopic = "users_created"

func (*Server) RegisterUser(context context.Context, request *grpc.RegisterUserRequest) (*grpc.UserProfileResponse, error) {
	conn := postgres.GetPostgresConnection()

	usernameExists, err := usecase.DoesUsernameExist(request.GetUsername(), conn)
	if err != nil {
		log.Printf("Error while checking username exists: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	if usernameExists {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}

	emailExists, err := usecase.DoesEmailExist(request.GetEmail(), conn)
	if err != nil {
		log.Printf("Error while checking email exists: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	if emailExists {
		return nil, status.Error(codes.AlreadyExists, "email already exists")
	}

	user, err := usecase.CreateUser(request.GetUsername(), request.GetEmail(), request.GetPassword(), conn)
	if err != nil {
		log.Printf("Error while creating user: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	registrationDateBytes := []byte(user.CreatedAt.Format("2006-01-02 15:04:05"))
	_, _, err = kafka.SendMessageSync(usersCreatedTopic, []byte(user.Id), registrationDateBytes, false)
	if err != nil {
		log.Printf("Error while sending message, reverting user creation: %v", err)
		err = usecase.DeleteUser(user.Id, conn)
		if err != nil {
			log.Printf("Error while reverting user creation: %v", err)
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &grpc.UserProfileResponse{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		LastLogin: timestamppb.New(user.LastLogin),
		IsActive:  true,
	}, nil
}
