package router

import (
	"vicomova/internal/interface/gateway/handler"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoutes(h *server.Hertz, userHandler *handler.UserHandler) {
	g := h.Group("/")
	users := g.Group("/user")
	{
		users.POST("/register", userHandler.Register)
		users.POST("/login", userHandler.Login)
		users.POST("/refresh", userHandler.RefreshToken)
		users.POST("/logout", userHandler.Logout)
		users.GET("/:id", userHandler.GetUser)
	}
}
