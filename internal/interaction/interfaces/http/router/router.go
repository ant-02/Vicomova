package router

import (
	"vicomova/internal/interaction/interfaces/http/handler"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoutes(h *server.Hertz, interactionHandler *handler.InteractionHandler) {
	g := h.Group("/", rootMw()...)

	video := g.Group("/video")
	video.POST("/like", interactionHandler.LikeVideo)
	video.DELETE("/like", interactionHandler.UnlikeVideo)
	video.GET("/like/list", interactionHandler.ListLikes)

	video.POST("/favorite", interactionHandler.AddFavorite)
	video.DELETE("/favorite", interactionHandler.RemoveFavorite)
	video.GET("/favorite/list", interactionHandler.ListFavorites)

	video.POST("/comment", interactionHandler.Comment)
	video.DELETE("/comment", interactionHandler.DeleteComment)
	video.GET("/comment/list", interactionHandler.ListComments)

	comment := g.Group("/comment")
	comment.POST("/like", interactionHandler.LikeComment)
}
