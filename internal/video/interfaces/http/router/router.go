package router

import (
	"vicomova/internal/shared/pkg/middleware"
	"vicomova/internal/user/domain/service"
	"vicomova/internal/video/interfaces/http/handler"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoutes(h *server.Hertz, videoHandler *handler.VideoHandler, tokenSvc *service.TokenService) {
	video := h.Group("/video")

	// Public routes (no auth required)
	video.GET("/stream", videoHandler.GetVideoStream)
	video.GET("/list", videoHandler.ListByCategory)
	video.GET("/hot", videoHandler.ListHotVideos)
	video.GET("/cover", videoHandler.GetVideoCover)

	// Protected routes (auth required)
	videoAuth := video.Group("/", middleware.Auth(tokenSvc))
	videoAuth.POST("/publish", videoHandler.PublishVideo)
	videoAuth.GET("/list/published", videoHandler.GetPublishedList)
}
