package router

import (
	"vicomova/internal/user/interfaces/http/handler"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoutes(h *server.Hertz, userHandler *handler.UserHandler) {
	g := h.Group("/")
	users := g.Group("/user")
	{
		users.POST("/register/send", userHandler.SendVerificationCode)
		users.POST("/register/verify", userHandler.VerifyAndRegister)
		users.POST("/login", userHandler.Login)
		users.POST("/refresh", userHandler.RefreshToken)
		users.POST("/logout", userHandler.Logout)
		users.GET("/:id", userHandler.GetUser)
	}
}