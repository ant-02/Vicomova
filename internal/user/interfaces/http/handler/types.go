package handler

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

// @Description 发送验证码请求
type SendCodeRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// @Description 发送验证码响应
type SendCodeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// @Description 验证注册请求
type VerifyCodeRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Code     string `json:"code"`
}

// @Description 注册响应
type RegisterResponse struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
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