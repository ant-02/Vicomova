package handler

import (
	"context"

	rpc "vicomova/internal/user/interfaces/grpc"
	hertz "vicomova/pkg/hertz"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type UserHandler struct {
	userClient *rpc.UserClient
}

func NewUserHandler(userClient *rpc.UserClient) *UserHandler {
	return &UserHandler{userClient: userClient}
}

// @Summary 用户登录
// @Description 用户登录获取Token
// @Tags user
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /user/login [post]
func (h *UserHandler) Login(ctx context.Context, c *app.RequestContext) {
	var req LoginRequest

	if err := c.Bind(&req); err != nil {
		hlog.Errorf("Login: invalid request body: %v", err)
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	resp, err := h.userClient.Login(ctx, req.Username, req.Password)
	if err != nil {
		hlog.Errorf("Login: username=%s failed: %v", req.Username, err)
		c.JSON(401, hertz.Fail(401, "Login failed"))
		return
	}
	hlog.Infof("Login: username=%s success, userID=%d", req.Username, resp.UserId)

	c.JSON(200, hertz.Success(map[string]interface{}{
		"user_id":       resp.UserId,
		"username":      resp.Username,
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"expires_in":    resp.ExpiresIn,
	}))
}

// @Summary 获取用户信息
// @Description 根据用户ID获取用户信息
// @Tags user
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /user/{id} [get]
func (h *UserHandler) GetUser(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("id")
	if userID == 0 {
		hlog.Error("GetUser: missing user_id")
		c.JSON(400, hertz.Fail(400, "Missing user_id"))
		return
	}

	resp, err := h.userClient.GetUser(ctx, userID)
	if err != nil {
		hlog.Errorf("GetUser: userID=%d failed: %v", userID, err)
		c.JSON(404, hertz.Fail(404, "User not found"))
		return
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"user_id":    resp.UserId,
		"username":   resp.Username,
		"email":      resp.Email,
		"created_at": resp.CreatedAt,
	}))
}

// @Summary 刷新Token
// @Description 使用refresh_token获取新的access_token
// @Tags user
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "刷新Token请求"
// @Success 200 {object} RefreshTokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /user/refresh [post]
func (h *UserHandler) RefreshToken(ctx context.Context, c *app.RequestContext) {
	var req RefreshTokenRequest

	if err := c.Bind(&req); err != nil {
		hlog.Errorf("RefreshToken: invalid request body: %v", err)
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	resp, err := h.userClient.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		hlog.Errorf("RefreshToken: failed: %v", err)
		c.JSON(401, hertz.Fail(401, "Refresh token failed"))
		return
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"expires_in":    resp.ExpiresIn,
	}))
}

// @Summary 用户登出
// @Description 登出并使refresh_token失效
// @Tags user
// @Accept json
// @Produce json
// @Param request body LogoutRequest true "登出请求"
// @Success 200 {object} LogoutResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /user/logout [post]
func (h *UserHandler) Logout(ctx context.Context, c *app.RequestContext) {
	var req LogoutRequest

	if err := c.Bind(&req); err != nil {
		hlog.Errorf("Logout: invalid request body: %v", err)
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	resp, err := h.userClient.Logout(ctx, req.AccessToken)
	if err != nil {
		hlog.Errorf("Logout: failed: %v", err)
		c.JSON(401, hertz.Fail(401, "Logout failed"))
		return
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"success": resp.Success,
	}))
}

// @Summary 发送邮箱验证码
// @Description 发送注册验证码到用户邮箱
// @Tags user
// @Accept json
// @Produce json
// @Param request body SendCodeRequest true "发送验证码请求"
// @Success 200 {object} SendCodeResponse
// @Failure 400 {object} ErrorResponse
// @Router /user/register/send [post]
func (h *UserHandler) SendVerificationCode(ctx context.Context, c *app.RequestContext) {
	var req SendCodeRequest

	if err := c.Bind(&req); err != nil {
		hlog.Errorf("SendVerificationCode: invalid request body: %v", err)
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	resp, err := h.userClient.SendVerificationCode(ctx, req.Username, req.Password, req.Email)
	if err != nil {
		hlog.Errorf("SendVerificationCode: username=%s email=%s failed: %v", req.Username, req.Email, err)
		c.JSON(500, hertz.Fail(500, "Send verification code failed"))
		return
	}
	hlog.Infof("SendVerificationCode: username=%s email=%s success=%v", req.Username, req.Email, resp.Success)

	c.JSON(200, hertz.Success(map[string]interface{}{
		"success": resp.Success,
		"message": resp.Message,
	}))
}

// @Summary 验证邮箱并完成注册
// @Description 使用邮箱验证码完成用户注册
// @Tags user
// @Accept json
// @Produce json
// @Param request body VerifyCodeRequest true "验证注册请求"
// @Success 200 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Router /user/register/verify [post]
func (h *UserHandler) VerifyAndRegister(ctx context.Context, c *app.RequestContext) {
	var req VerifyCodeRequest

	if err := c.Bind(&req); err != nil {
		hlog.Errorf("VerifyAndRegister: invalid request body: %v", err)
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	resp, err := h.userClient.VerifyAndRegister(ctx, req.Username, req.Password, req.Email, req.Code)
	if err != nil {
		hlog.Errorf("VerifyAndRegister: username=%s email=%s failed: %v", req.Username, req.Email, err)
		c.JSON(500, hertz.Fail(500, "Verify and register failed"))
		return
	}
	hlog.Infof("VerifyAndRegister: username=%s success, userID=%d", req.Username, resp.UserId)

	c.JSON(200, hertz.Success(map[string]interface{}{
		"user_id":  resp.UserId,
		"username": resp.Username,
	}))
}