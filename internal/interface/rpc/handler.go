package rpc

import (
	"context"

	user "vicomova/third_party/kitex_gen/user"
	appuser "vicomova/internal/application/user"
)

type UserHandler struct {
	cmdSvc *appuser.UserCommandService
	querySvc *appuser.UserQueryService
}

func NewUserHandler(cmdSvc *appuser.UserCommandService, querySvc *appuser.UserQueryService) *UserHandler {
	return &UserHandler{
		cmdSvc:  cmdSvc,
		querySvc: querySvc,
	}
}

func (h *UserHandler) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	cmd := &appuser.RegisterCommand{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	result, err := h.cmdSvc.Register(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &user.RegisterResponse{
		UserId:   result.UserID,
		Username: result.Username,
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	cmd := &appuser.LoginCommand{
		Username: req.Username,
		Password: req.Password,
	}

	result, err := h.cmdSvc.Login(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &user.LoginResponse{
		UserId:       result.UserID,
		Username:     result.Username,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	query := &appuser.GetUserQuery{
		UserID: req.UserId,
	}

	result, err := h.querySvc.GetUser(ctx, query)
	if err != nil {
		return nil, err
	}

	return &user.GetUserResponse{
		UserId:   result.UserID,
		Username: result.Username,
		Email:    result.Email,
	}, nil
}

func (h *UserHandler) RefreshToken(ctx context.Context, req *user.RefreshTokenRequest) (*user.RefreshTokenResponse, error) {
	cmd := &appuser.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
	}

	result, err := h.cmdSvc.RefreshToken(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &user.RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (h *UserHandler) Logout(ctx context.Context, req *user.LogoutRequest) (*user.LogoutResponse, error) {
	cmd := &appuser.LogoutCommand{
		AccessToken: req.AccessToken,
	}

	err := h.cmdSvc.Logout(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &user.LogoutResponse{
		Success: true,
	}, nil
}
