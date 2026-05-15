package http

import (
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoutes(h *server.Hertz, searchHandler *SearchHandler) {
	search := h.Group("/search")

	search.GET("", searchHandler.SearchVideos)
	search.POST("/index", searchHandler.IndexVideo)
	search.DELETE("/index/:video_id", searchHandler.DeleteVideoIndex)
}
