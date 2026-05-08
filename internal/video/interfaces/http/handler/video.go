package handler

import (
	"context"
	"strconv"

	videoRpc "vicomova/internal/video/interfaces/grpc"
	hertz "vicomova/pkg/infrastructure/hertz"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type VideoHandler struct {
	videoClient *videoRpc.VideoClient
}

func NewVideoHandler(videoClient *videoRpc.VideoClient) *VideoHandler {
	return &VideoHandler{videoClient: videoClient}
}

// @Summary 获取视频流
// @Description 获取视频播放地址
// @Tags video
// @Produce json
// @Param video_id query int64 true "视频ID"
// @Success 200 {object} VideoStreamResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/stream [get]
func (h *VideoHandler) GetVideoStream(ctx context.Context, c *app.RequestContext) {
	videoID, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid video_id"))
		return
	}

	// Increment view count
	go func() {
		_, _ = h.videoClient.IncrementView(context.Background(), videoID)
	}()

	resp, err := h.videoClient.GetVideoStream(ctx, videoID)
	if err != nil {
		hlog.Errorf("GetVideoStream: videoID=%d failed: %v", videoID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get video stream"))
		return
	}

	c.JSON(200, hertz.Success(VideoStreamResponse{
		VideoURL: resp.VideoUrl,
		Title:    resp.Title,
	}))
}

// @Summary 按类别获取视频列表
// @Description 分页获取指定类别的视频列表
// @Tags video
// @Produce json
// @Param category_id query int32 true "分类ID"
// @Param page query int32 false "页码" default(1)
// @Param size query int32 false "每页数量" default(10)
// @Success 200 {object} VideoListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/list [get]
func (h *VideoHandler) ListByCategory(ctx context.Context, c *app.RequestContext) {
	categoryID, err := strconv.ParseInt(c.Query("category_id"), 10, 64)
	if err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid category_id"))
		return
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	size, _ := strconv.ParseInt(c.DefaultQuery("size", "10"), 10, 32)

	resp, err := h.videoClient.ListByCategory(ctx, int32(categoryID), int32(page), int32(size))
	if err != nil {
		hlog.Errorf("ListByCategory: categoryID=%d failed: %v", categoryID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get video list"))
		return
	}

	videos := make([]*VideoItem, 0, len(resp.Videos))
	for _, v := range resp.Videos {
		videos = append(videos, &VideoItem{
			ID:           v.Id,
			UserID:       v.UserId,
			Title:        v.Title,
			Description:  v.Description,
			CoverURL:     v.CoverUrl,
			VideoURL:     v.VideoUrl,
			CategoryID:   v.CategoryId,
			ViewCount:    v.ViewCount,
			LikeCount:    v.LikeCount,
			CommentCount: v.CommentCount,
			Duration:     v.Duration,
			CreatedAt:    v.CreatedAt,
		})
	}

	c.JSON(200, hertz.Success(VideoListResponse{
		Videos: videos,
		Total:  resp.Total,
	}))
}

// @Summary 获取热门视频
// @Description 获取热门视频列表（基于热度算法）
// @Tags video
// @Produce json
// @Param limit query int32 false "返回数量" default(10)
// @Success 200 {object} VideoListResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/hot [get]
func (h *VideoHandler) ListHotVideos(ctx context.Context, c *app.RequestContext) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 32)

	resp, err := h.videoClient.ListHotVideos(ctx, int32(limit))
	if err != nil {
		hlog.Errorf("ListHotVideos: failed: %v", err)
		c.JSON(500, hertz.Fail(500, "Failed to get hot videos"))
		return
	}

	videos := make([]*VideoItem, 0, len(resp.Videos))
	for _, v := range resp.Videos {
		videos = append(videos, &VideoItem{
			ID:           v.Id,
			UserID:       v.UserId,
			Title:        v.Title,
			Description:  v.Description,
			CoverURL:     v.CoverUrl,
			VideoURL:     v.VideoUrl,
			CategoryID:   v.CategoryId,
			ViewCount:    v.ViewCount,
			LikeCount:    v.LikeCount,
			CommentCount: v.CommentCount,
			Duration:     v.Duration,
			CreatedAt:    v.CreatedAt,
		})
	}

	c.JSON(200, hertz.Success(VideoListResponse{
		Videos: videos,
		Total:  int64(len(videos)),
	}))
}

// @Summary 投稿视频
// @Description 发布新视频（需要认证）
// @Tags video
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body PublishVideoRequest true "投稿信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/publish [post]
func (h *VideoHandler) PublishVideo(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req PublishVideoRequest
	if err := c.Bind(&req); err != nil {
		hlog.Errorf("PublishVideo: invalid request: %v", err)
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	resp, err := h.videoClient.PublishVideo(ctx, userID, req.Title, req.Description, req.CategoryId, req.CoverUrl, req.VideoUrl, req.Duration)
	if err != nil {
		hlog.Errorf("PublishVideo: userID=%d failed: %v", userID, err)
		c.JSON(500, hertz.Fail(500, "Failed to publish video"))
		return
	}

	hlog.Infof("PublishVideo: userID=%d success, videoID=%d", userID, resp.VideoId)
	c.JSON(200, hertz.Success(map[string]interface{}{
		"video_id": resp.VideoId,
		"success":  true,
	}))
}

// @Summary 获取视频封面
// @Description 获取视频封面地址
// @Tags video
// @Produce json
// @Param video_id query int64 true "视频ID"
// @Success 200 {object} CoverResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/cover [get]
func (h *VideoHandler) GetVideoCover(ctx context.Context, c *app.RequestContext) {
	videoID, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid video_id"))
		return
	}

	resp, err := h.videoClient.GetVideoCover(ctx, videoID)
	if err != nil {
		hlog.Errorf("GetVideoCover: videoID=%d failed: %v", videoID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get video cover"))
		return
	}

	c.JSON(200, hertz.Success(CoverResponse{
		CoverURL: resp.CoverUrl,
	}))
}

// @Summary 获取发布列表
// @Description 获取用户发布的视频列表（需要认证）
// @Tags video
// @Produce json
// @Security BearerAuth
// @Param page query int32 false "页码" default(1)
// @Param size query int32 false "每页数量" default(10)
// @Success 200 {object} VideoListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/list/published [get]
func (h *VideoHandler) GetPublishedList(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	size, _ := strconv.ParseInt(c.DefaultQuery("size", "10"), 10, 32)

	resp, err := h.videoClient.GetPublishedList(ctx, userID, int32(page), int32(size))
	if err != nil {
		hlog.Errorf("GetPublishedList: userID=%d failed: %v", userID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get published list"))
		return
	}

	videos := make([]*VideoItem, 0, len(resp.Videos))
	for _, v := range resp.Videos {
		videos = append(videos, &VideoItem{
			ID:           v.Id,
			UserID:       v.UserId,
			Title:        v.Title,
			Description:  v.Description,
			CoverURL:     v.CoverUrl,
			VideoURL:     v.VideoUrl,
			CategoryID:   v.CategoryId,
			ViewCount:    v.ViewCount,
			LikeCount:    v.LikeCount,
			CommentCount: v.CommentCount,
			Duration:     v.Duration,
			CreatedAt:    v.CreatedAt,
		})
	}

	c.JSON(200, hertz.Success(VideoListResponse{
		Videos: videos,
		Total:  resp.Total,
	}))
}
