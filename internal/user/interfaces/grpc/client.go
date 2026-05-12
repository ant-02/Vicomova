package rpc

import (
	"context"

	user "vicomova/third_party/kitex_gen/user"
	userservice "vicomova/third_party/kitex_gen/user/userservice"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
)

type UserClient struct {
	cli userservice.Client
}

// NewUserClient 创建 RPC 客户端，直连指定地址（用于本地开发无服务注册中心的场景）
func NewUserClient(serviceName, addr string) (*UserClient, error) {
	cli, err := userservice.NewClient(serviceName,
		client.WithHostPorts(addr),
	)
	if err != nil {
		return nil, err
	}

	return &UserClient{cli: cli}, nil
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

func (c *UserClient) RefreshToken(ctx context.Context, refreshToken string) (*user.RefreshTokenResponse, error) {
	return c.cli.RefreshToken(ctx, &user.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
}

func (c *UserClient) Logout(ctx context.Context, refreshToken string) (*user.LogoutResponse, error) {
	klog.Debugf("UserClient.Logout: refreshToken=%s", refreshToken)
	return c.cli.Logout(ctx, &user.LogoutRequest{
		RefreshToken: refreshToken,
	})
}

func (c *UserClient) SendVerificationCode(ctx context.Context, username, password, email string) (*user.SendVerificationCodeResponse, error) {
	return c.cli.SendVerificationCode(ctx, &user.SendVerificationCodeRequest{
		Username: username,
		Password: password,
		Email:    email,
	})
}

func (c *UserClient) VerifyAndRegister(ctx context.Context, username, password, email, code string) (*user.RegisterResponse, error) {
	return c.cli.VerifyAndRegister(ctx, &user.VerifyAndRegisterRequest{
		Username: username,
		Password: password,
		Email:    email,
		Code:     code,
	})
}

func (c *UserClient) BatchGetUsers(ctx context.Context, userIDs []int64) (map[int64]*user.User, error) {
	resp, err := c.cli.BatchGetUsers(ctx, &user.BatchGetUsersRequest{
		UserIds: userIDs,
	})
	if err != nil {
		return nil, err
	}
	result := make(map[int64]*user.User, len(resp.Users))
	for _, u := range resp.Users {
		result[u.UserId] = u
	}
	return result, nil
}
