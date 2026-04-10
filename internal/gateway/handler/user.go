package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	hertz "vicomova/pkg/hertz"
	"vicomova/internal/rpc/user"
)

type UserHandler struct {
	userClient *user.UserClient
}

func NewUserHandler(userClient *user.UserClient) *UserHandler {
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
		"user_id":  resp.UserId,
		"username": resp.Username,
		"token":    resp.Token,
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
		"user_id":   resp.UserId,
		"username":  resp.Username,
		"email":     resp.Email,
		"created_at": resp.CreatedAt,
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
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Token    string `json:"token"`
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
