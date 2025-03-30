package usecase

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	gen "sn/gateway/generated"
	"sn/libraries/proto/users"
	"time"
)

// Requires defer cl() to be called after
func getUsersClient() (users.UserServiceClient, context.Context, func(), error) {
	c, err := grpc.NewClient("users:50002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, func() {}, err
	}

	client := users.NewUserServiceClient(c)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	return client, ctx, func() { c.Close(); cancel() }, nil
}

func RegisterUser(body *gen.UserRegistration) (*users.UserProfileResponse, error) {
	client, ctx, cl, err := getUsersClient()
	if err != nil {
		return nil, err
	}
	defer cl()

	response, err := client.RegisterUser(ctx, &users.RegisterUserRequest{
		Username: body.Username,
		Password: body.Password,
		Email:    string(body.Email),
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func AuthUser(body gen.UsernameAndPassword) (string, error) {
	client, ctx, cl, err := getUsersClient()
	if err != nil {
		return "", err
	}
	defer cl()

	response, err := client.AuthenticateUser(ctx, &users.AuthenticateUserRequest{
		Username: body.Username,
		Password: body.Password,
	})
	if err != nil {
		return "", err
	}

	return response.GetId(), nil
}

func fillProfile(userId, email, username string,
	lastLogin *timestamppb.Timestamp, dateOfBirth, fullName, phoneNumber string) (*gen.ProfileResponse, error) {
	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	result := &gen.ProfileResponse{
		Id:       parsedUserId,
		Email:    email,
		Username: username,
	}
	if lastLogin != nil {
		result.LastLogin = lastLogin.AsTime().Format("2006-01-02 15:04")
	}
	if dateOfBirth != "" {
		parsedDate, err := time.Parse(time.RFC3339, dateOfBirth)
		if err != nil {
			return nil, err
		}
		parsedDateStr := parsedDate.Format("2006-01-02")
		result.DateOfBirth = &parsedDateStr
	}
	if fullName != "" {
		result.FullName = &fullName
	}
	if phoneNumber != "" {
		result.PhoneNumber = &phoneNumber
	}

	return result, nil
}

func GetProfile(userId string) (*gen.ProfileResponse, error) {
	client, ctx, cl, err := getUsersClient()
	if err != nil {
		return nil, err
	}
	defer cl()

	response, err := client.GetUserProfile(ctx, &users.GetUserProfileRequest{
		Id: userId,
	})
	if err != nil {
		return nil, err
	}

	return fillProfile(userId, response.GetEmail(), response.GetUsername(), response.GetLastLogin(), response.GetDateOfBirth(),
		response.GetFullName(), response.GetPhoneNumber())
}

func EditProfile(userId string, body *gen.EditProfile) (*gen.ProfileResponse, error) {
	client, ctx, cl, err := getUsersClient()
	if err != nil {
		return nil, err
	}
	defer cl()

	request := users.UpdateUserProfileRequest{
		Id: userId,
	}
	if body.DateOfBirth != nil {
		request.DateOfBirth = body.DateOfBirth.Format("2006-01-02")
	}
	if body.FullName != nil {
		request.FullName = *body.FullName
	}
	if body.PhoneNumber != nil {
		request.PhoneNumber = *body.PhoneNumber
	}

	response, err := client.UpdateUserProfile(ctx, &request)
	if err != nil {
		return nil, err
	}

	return fillProfile(userId, response.GetEmail(), response.GetUsername(), response.GetLastLogin(), response.GetDateOfBirth(),
		response.GetFullName(), response.GetPhoneNumber())
}
