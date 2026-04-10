package user

import (
	"context"

	user "vicomova/third_party/kitex_gen/user"
	userservice "vicomova/third_party/kitex_gen/user/userservice"
)

type UserClient struct {
	cli userservice.Client
}

func NewUserClient(addr string) (*UserClient, error) {
	cli, err := userservice.NewClient(
		"user",
	)
	if err != nil {
		return nil, err
	}

	return &UserClient{cli: cli}, nil
}

func (c *UserClient) Register(ctx context.Context, username, password, email string) (*user.RegisterResponse, error) {
	return c.cli.Register(ctx, &user.RegisterRequest{
		Username: username,
		Password: password,
		Email:    email,
	})
}

func (c *UserClient) Login(ctx context.Context, username, password string) (*user.LoginResponse, error) {
	return c.cli.Login(ctx, &user.LoginRequest{
		Username: username,
		Password: password,
	})
}

func (c *UserClient) GetUser(ctx context.Context, userID int64) (*user.GetUserResponse, error) {
	return c.cli.GetUser(ctx, &user.GetUserRequest{
		UserId: userID,
	})
}
