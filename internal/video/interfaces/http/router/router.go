package router

import (
	"vicomova/internal/video/interfaces/http/handler"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoutes(h *server.Hertz, videoHandler *handler.VideoHandler) {
	video := h.Group("/video")

	// Public routes (no auth required)
	video.GET("/stream", videoHandler.GetVideoStream)
	video.GET("/hot", videoHandler.ListHotVideos)
	video.GET("/list/published", videoHandler.GetPublishedList)
	video.GET("/category/:id/list", videoHandler.ListCategoryVideos)

	// Protected routes (auth required)
	videoAuth := video.Group("/", videoAuthMw()...)
	videoAuth.POST("/save", videoHandler.SaveVideo)
	videoAuth.POST("/submit", videoHandler.SubmitVideo)
	videoAuth.POST("/publish", videoHandler.PublishVideo)
	videoAuth.POST("/upload/token", videoHandler.GetUploadToken)
}
