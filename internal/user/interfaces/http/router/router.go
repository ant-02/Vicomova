package router

import (
	"vicomova/internal/user/interfaces/http/handler"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoutes(h *server.Hertz, userHandler *handler.UserHandler) {
	users := h.Group("/user")

	// 公开路由 - 无需认证
	public := users.Group("/")
	public.POST("/register/send", userHandler.SendVerificationCode)
	public.POST("/register/verify", userHandler.VerifyAndRegister)
	public.POST("/login", userHandler.Login)
	public.POST("/refresh", userHandler.RefreshToken)

	// 受保护路由 - 需要认证
	protected := users.Group("/", protectedMw()...)
	protected.GET("", userHandler.GetUser)
	protected.POST("/logout", userHandler.Logout)
}
