package http

import (
	"context"
	"strconv"

	"vicomova/internal/search/application/query"
	"vicomova/pkg/log"

	"github.com/cloudwego/hertz/pkg/app"
)

type SearchHandler struct {
	svc *query.SearchService
}

func NewSearchHandler(svc *query.SearchService) *SearchHandler {
	return &SearchHandler{svc: svc}
}

func (h *SearchHandler) SearchVideos(ctx context.Context, c *app.RequestContext) {
	queryStr := c.Query("q")
	limitStr := c.DefaultQuery("limit", "20")
	cursor := c.Query("cursor")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	q := &query.SearchVideosQuery{
		Query:  queryStr,
		Limit:  limit,
		Cursor: cursor,
	}

	result, err := h.svc.SearchVideos(ctx, q)
	if err != nil {
		log.Error.Printf("SearchVideos failed: %v", err)
		c.JSON(500, map[string]any{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

func (h *SearchHandler) IndexVideo(ctx context.Context, c *app.RequestContext) {
	var req struct {
		VideoID     int64  `json:"video_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]any{"error": err.Error()})
		return
	}

	q := &query.IndexVideoQuery{
		VideoID:     req.VideoID,
		Title:       req.Title,
		Description: req.Description,
	}

	if err := h.svc.IndexVideo(ctx, q); err != nil {
		log.Error.Printf("IndexVideo failed: %v", err)
		c.JSON(500, map[string]any{"error": err.Error()})
		return
	}

	c.JSON(200, map[string]any{"status": "ok"})
}

func (h *SearchHandler) DeleteVideoIndex(ctx context.Context, c *app.RequestContext) {
	videoIDStr := c.Param("video_id")
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(400, map[string]any{"error": "invalid video_id"})
		return
	}

	if err := h.svc.DeleteVideoIndex(ctx, videoID); err != nil {
		log.Error.Printf("DeleteVideoIndex failed: %v", err)
		c.JSON(500, map[string]any{"error": err.Error()})
		return
	}

	c.JSON(200, map[string]any{"status": "ok"})
}
