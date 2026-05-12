package rpc

import (
	"context"

	appCommand "vicomova/internal/user/application/command"
	appQuery "vicomova/internal/user/application/query"
	user "vicomova/third_party/kitex_gen/user"

	"github.com/cloudwego/kitex/pkg/klog"
)

type UserHandler struct {
	cmdSvc   *appCommand.UserCommandService
	querySvc *appQuery.UserQueryService
}

func NewUserHandler(cmdSvc *appCommand.UserCommandService, querySvc *appQuery.UserQueryService) *UserHandler {
	return &UserHandler{
		cmdSvc:   cmdSvc,
		querySvc: querySvc,
	}
}

func (h *UserHandler) SendVerificationCode(ctx context.Context, req *user.SendVerificationCodeRequest) (*user.SendVerificationCodeResponse, error) {
	cmd := &appCommand.SendVerificationCodeCommand{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	result, err := h.cmdSvc.SendVerificationCode(ctx, cmd)
	if err != nil {
		klog.Errorf("SendVerificationCode: username=%s email=%s failed: %v", req.Username, req.Email, err)
		return nil, err
	}
	klog.Infof("SendVerificationCode: username=%s email=%s success=%v", req.Username, req.Email, result.Success)

	return &user.SendVerificationCodeResponse{
		Success: result.Success,
		Message: result.Message,
	}, nil
}

func (h *UserHandler) VerifyAndRegister(ctx context.Context, req *user.VerifyAndRegisterRequest) (*user.RegisterResponse, error) {
	cmd := &appCommand.VerifyAndRegisterCommand{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Code:     req.Code,
	}

	result, err := h.cmdSvc.VerifyAndRegister(ctx, cmd)
	if err != nil {
		klog.Errorf("VerifyAndRegister: username=%s email=%s failed: %v", req.Username, req.Email, err)
		return nil, err
	}
	klog.Infof("VerifyAndRegister: username=%s success, userID=%d", req.Username, result.UserID)

	return &user.RegisterResponse{
		UserId:   result.UserID,
		Username: result.Username,
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	cmd := &appCommand.LoginCommand{
		Username: req.Username,
		Password: req.Password,
	}

	result, err := h.cmdSvc.Login(ctx, cmd)
	if err != nil {
		klog.Errorf("Login: username=%s failed: %v", req.Username, err)
		return nil, err
	}
	klog.Infof("Login: username=%s success, userID=%d", req.Username, result.UserID)

	return &user.LoginResponse{
		UserId:       result.UserID,
		Username:     result.Username,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	query := &appQuery.GetUserQuery{
		UserID: req.UserId,
	}

	result, err := h.querySvc.GetUser(ctx, query)
	if err != nil {
		klog.Errorf("GetUser: userID=%d failed: %v", req.UserId, err)
		return nil, err
	}

	return &user.GetUserResponse{
		UserId:   result.UserID,
		Username: result.Username,
		Email:    result.Email,
	}, nil
}

func (h *UserHandler) RefreshToken(ctx context.Context, req *user.RefreshTokenRequest) (*user.RefreshTokenResponse, error) {
	cmd := &appCommand.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
	}

	result, err := h.cmdSvc.RefreshToken(ctx, cmd)
	if err != nil {
		klog.Errorf("RefreshToken: failed: %v", err)
		return nil, err
	}

	return &user.RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

func (h *UserHandler) Logout(ctx context.Context, req *user.LogoutRequest) (*user.LogoutResponse, error) {
	klog.Debugf("UserHandler.Logout: RefreshToken=%s", req.RefreshToken)
	cmd := &appCommand.LogoutCommand{
		RefreshToken: req.RefreshToken,
	}

	err := h.cmdSvc.Logout(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &user.LogoutResponse{
		Success: true,
	}, nil
}

func (h *UserHandler) BatchGetUsers(ctx context.Context, req *user.BatchGetUsersRequest) (*user.BatchGetUsersResponse, error) {
	query := &appQuery.BatchGetUsersQuery{
		UserIDs: req.UserIds,
	}
	results, err := h.querySvc.BatchGetUsers(ctx, query)
	if err != nil {
		klog.Errorf("BatchGetUsers: failed: %v", err)
		return nil, err
	}
	users := make([]*user.User, 0, len(results))
	for _, r := range results {
		users = append(users, &user.User{
			UserId:   r.UserID,
			Username: r.Username,
			Email:    r.Email,
		})
	}
	return &user.BatchGetUsersResponse{Users: users}, nil
}
