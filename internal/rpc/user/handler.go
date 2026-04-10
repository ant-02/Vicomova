package user

import (
	"context"

	user "vicomova/third_party/kitex_gen/user"
	"vicomova/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	input := &service.RegisterInput{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	result, err := h.svc.Register(ctx, input)
	if err != nil {
		return nil, err
	}

	return &user.RegisterResponse{
		UserId:   result.UserID,
		Username: result.Username,
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	input := &service.LoginInput{
		Username: req.Username,
		Password: req.Password,
	}

	result, err := h.svc.Login(ctx, input)
	if err != nil {
		return nil, err
	}

	return &user.LoginResponse{
		UserId:   result.UserID,
		Username: result.Username,
		Token:    result.Token,
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	result, err := h.svc.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &user.GetUserResponse{
		UserId:   result.UserID,
		Username: result.Username,
		Email:    result.Email,
	}, nil
}
