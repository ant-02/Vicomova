package router

import (
	"vicomova/internal/gateway/handler"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoutes(h *server.Hertz, userHandler *handler.UserHandler) {
	g := h.Group("/")
	users := g.Group("/user")
	{
		users.POST("/register", userHandler.Register)
		users.POST("/login", userHandler.Login)
		users.GET("/:id", userHandler.GetUser)
	}
}
