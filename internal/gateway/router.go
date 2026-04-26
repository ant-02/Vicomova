package gateway

import (
	"vicomova/internal/gateway/handlers"
	"vicomova/internal/shared/pkg/middleware"
	"vicomova/internal/user/domain/service"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterVideoAndInteractionRoutes(
	h *server.Hertz,
	videoHandler *handlers.VideoHandler,
	interactionHandler *handlers.InteractionHandler,
	tokenSvc *service.TokenService,
) {
	g := h.Group("/")

	// ========== Video Routes ==========
	video := g.Group("/video")

	// Public video routes (no auth required)
	video.GET("/stream", videoHandler.GetVideoStream)
	video.GET("/list", videoHandler.ListByCategory)
	video.GET("/hot", videoHandler.ListHotVideos)
	video.GET("/cover", videoHandler.GetVideoCover)

	// Protected video routes (auth required)
	videoAuth := video.Group("/", middleware.Auth(tokenSvc))
	videoAuth.POST("/publish", videoHandler.PublishVideo)
	videoAuth.GET("/list/published", videoHandler.GetPublishedList)

	// ========== Video Interaction Routes ==========
	// Like routes
	video.POST("/like", interactionHandler.LikeVideo)
	video.DELETE("/like", interactionHandler.UnlikeVideo)
	video.GET("/like/list", interactionHandler.ListLikes)

	// Favorite routes
	video.POST("/favorite", interactionHandler.AddFavorite)
	video.DELETE("/favorite", interactionHandler.RemoveFavorite)
	video.GET("/favorite/list", interactionHandler.ListFavorites)

	// Comment routes
	video.POST("/comment", interactionHandler.Comment)
	video.DELETE("/comment", interactionHandler.DeleteComment)
	video.GET("/comment/list", interactionHandler.ListComments)

	// ========== Comment Interaction Routes ==========
	comment := g.Group("/comment")
	comment.POST("/like", interactionHandler.LikeComment)
}
