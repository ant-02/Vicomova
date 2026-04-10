package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	hertz "vicomova/pkg/hertz"
	"vicomova/internal/interface/rpc"
)

type UserHandler struct {
	userClient *rpc.UserClient
}

func NewUserHandler(userClient *rpc.UserClient) *UserHandler {
	return &UserHandler{userClient: userClient}
}

// @Summary 用户注册
// @Description 创建新用户
// @Tags user
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 200 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /user/register [post]
func (h *UserHandler) Register(ctx context.Context, c *app.RequestContext) {
	var req RegisterRequest

	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	resp, err := h.userClient.Register(ctx, req.Username, req.Password, req.Email)
	if err != nil {
		c.JSON(500, hertz.Fail(500, "Register failed"))
		return
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"user_id":  resp.UserId,
		"username": resp.Username,
	}))
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
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	resp, err := h.userClient.Login(ctx, req.Username, req.Password)
	if err != nil {
		c.JSON(401, hertz.Fail(401, "Login failed"))
		return
	}

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
		c.JSON(400, hertz.Fail(400, "Missing user_id"))
		return
	}

	resp, err := h.userClient.GetUser(ctx, userID)
	if err != nil {
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
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	resp, err := h.userClient.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
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
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	resp, err := h.userClient.Logout(ctx, req.AccessToken)
	if err != nil {
		c.JSON(401, hertz.Fail(401, "Logout failed"))
		return
	}

	c.JSON(200, hertz.Success(map[string]interface{}{
		"success": resp.Success,
	}))
}

// @Description 注册请求
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// @Description 注册响应
type RegisterResponse struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

// @Description 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Description 登录响应
type LoginResponse struct {
	UserID       int64  `json:"user_id"`
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// @Description 刷新Token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// @Description 刷新Token响应
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// @Description 登出请求
type LogoutRequest struct {
	AccessToken string `json:"access_token"`
}

// @Description 登出响应
type LogoutResponse struct {
	Success bool `json:"success"`
}

// @Description 用户响应
type UserResponse struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
}

// @Description 错误响应
type ErrorResponse struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
}
