package sn_go

//go:generate mockgen -destination=libraries/proto/posts/mocks/posts_client_mock.go -package=posts_mocks sn/libraries/proto/posts PostServiceClient
//go:generate mockgen -destination=libraries/proto/users/mocks/users_client_mock.go -package=users_mocks sn/libraries/proto/users UserServiceClient
