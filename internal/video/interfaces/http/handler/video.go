package handler

import (
	"context"
	"strconv"
	"strings"

	videoRpc "vicomova/internal/video/interfaces/grpc"
	"vicomova/pkg/constants"
	hertz "vicomova/pkg/infrastructure/hertz"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// BizError 简化版，用于解析 gRPC 返回的错误
type BizError struct {
	Code int32
	Msg  string
}

// parseBizError 从 gRPC 错误中解析 BizError
// gRPC 错误格式: "biz error: code: 403, msg: Video not published"
func parseBizError(err error) *BizError {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if !strings.Contains(msg, "biz error:") {
		return nil
	}

	// 提取 code
	codeIdx := strings.Index(msg, "code: ")
	if codeIdx == -1 {
		return nil
	}
	codeStr := msg[codeIdx+6:]
	code, _ := strconv.ParseInt(codeStr[:strings.Index(codeStr, ",")], 10, 32)

	// 提取 msg
	msgIdx := strings.Index(msg, "msg: ")
	if msgIdx == -1 {
		return nil
	}
	msgStr := msg[msgIdx+5:]

	return &BizError{Code: int32(code), Msg: msgStr}
}

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

	resp, err := h.videoClient.GetVideoStream(ctx, videoID)
	if err != nil {
		hlog.Errorf("GetVideoStream: videoID=%d failed: %v", videoID, err)
		// 尝试解析 BizError
		if bizErr := parseBizError(err); bizErr != nil {
			c.JSON(int(bizErr.Code), hertz.Fail(bizErr.Code, bizErr.Msg))
			return
		}
		c.JSON(500, hertz.Fail(500, "Failed to get video stream"))
		return
	}

	c.JSON(200, hertz.Success(VideoStreamResponse{
		ID:           resp.Id,
		UserID:       resp.UserId,
		Title:        resp.Title,
		Description:  resp.Description,
		CoverURL:     resp.CoverUrl,
		VideoURL:     resp.VideoUrl,
		CategoryID:   resp.CategoryId,
		ViewCount:    resp.ViewCount,
		LikeCount:    resp.LikeCount,
		CommentCount: resp.CommentCount,
		Duration:     resp.Duration,
	}))
}

// @Summary 获取分类视频
// @Description 获取指定分类的视频列表，支持游标分页
// @Tags video
// @Produce json
// @Param id path int true "分类ID"
// @Param limit query int32 false "返回数量" default(10)
// @Param cursor query string false "游标分页（首次不传）" default("")
// @Success 200 {object} SuccessResponse{data=HotVideoListResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/category/{id}/list [get]
func (h *VideoHandler) ListCategoryVideos(ctx context.Context, c *app.RequestContext) {
	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil || categoryID <= 0 {
		c.JSON(400, hertz.Fail(400, "Invalid category_id"))
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 32)
	cursor := c.Query("cursor")

	resp, err := h.videoClient.ListCategoryVideos(ctx, int32(categoryID), int32(limit), cursor)
	if err != nil {
		hlog.Errorf("ListCategoryVideos: categoryID=%d failed: %v", categoryID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get category videos"))
		return
	}

	videos := make([]*HotVideoItem, 0, len(resp.Videos))
	for _, v := range resp.Videos {
		videos = append(videos, &HotVideoItem{
			ID:           v.Id,
			Title:        v.Title,
			CoverURL:     v.CoverUrl,
			Duration:     int(v.Duration),
			ViewCount:    v.ViewCount,
			CommentCount: v.CommentCount,
		})
	}

	c.JSON(200, hertz.Success(HotVideoListResponse{
		Videos:     videos,
		NextCursor: resp.NextCursor,
		HasMore:    resp.HasMore,
	}))
}

// @Summary 获取热门视频
// @Description 获取热门视频列表（基于热度算法），支持游标分页
// @Tags video
// @Produce json
// @Param limit query int32 false "返回数量" default(10)
// @Param cursor query string false "游标分页（首次不传）" default("")
// @Success 200 {object} SuccessResponse{data=HotVideoListResponse}
// @Failure 500 {object} ErrorResponse
// @Router /video/hot [get]
func (h *VideoHandler) ListHotVideos(ctx context.Context, c *app.RequestContext) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 32)
	cursor := c.Query("cursor")

	resp, err := h.videoClient.ListHotVideos(ctx, int32(limit), cursor)
	if err != nil {
		hlog.Errorf("ListHotVideos: failed: %v", err)
		c.JSON(500, hertz.Fail(500, "Failed to get hot videos"))
		return
	}

	videos := make([]*HotVideoItem, 0, len(resp.Videos))
	for _, v := range resp.Videos {
		videos = append(videos, &HotVideoItem{
			ID:           v.Id,
			Title:        v.Title,
			CoverURL:     v.CoverUrl,
			Duration:     int(v.Duration),
			ViewCount:    v.ViewCount,
			CommentCount: v.CommentCount,
			UserName:     v.UserName,
		})
	}

	c.JSON(200, hertz.Success(HotVideoListResponse{
		Videos:     videos,
		NextCursor: resp.NextCursor,
		HasMore:    resp.HasMore,
	}))
}

// @Summary 保存视频草稿
// @Description 创建或更新视频草稿，状态为编辑中（需要认证）
// @Tags video
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SaveVideoRequest true "保存信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/save [post]
func (h *VideoHandler) SaveVideo(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req SaveVideoRequest
	if err := c.Bind(&req); err != nil {
		hlog.Errorf("SaveVideo: invalid request: %v", err)
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	resp, err := h.videoClient.SaveVideo(ctx, req.VideoID, userID, req.Title, req.Description, req.CoverUrl, req.VideoUrl, req.Duration, req.CategoryId)
	if err != nil {
		hlog.Errorf("SaveVideo: userID=%d failed: %v", userID, err)
		c.JSON(500, hertz.Fail(500, "Failed to save video"))
		return
	}

	hlog.Infof("SaveVideo: userID=%d success, videoID=%d", userID, resp.VideoId)
	c.JSON(200, hertz.Success(map[string]interface{}{
		"video_id": resp.VideoId,
		"success":  true,
	}))
}

// @Summary 提交审核
// @Description 将草稿提交审核，状态从编辑中变为审核中（需要认证）
// @Tags video
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SubmitVideoRequest true "提交信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/submit [post]
func (h *VideoHandler) SubmitVideo(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req SubmitVideoRequest
	if err := c.Bind(&req); err != nil {
		hlog.Errorf("SubmitVideo: invalid request: %v", err)
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	resp, err := h.videoClient.SubmitVideo(ctx, req.VideoID, userID)
	if err != nil {
		hlog.Errorf("SubmitVideo: userID=%d videoID=%d failed: %v", userID, req.VideoID, err)
		c.JSON(500, hertz.Fail(500, "Failed to submit video"))
		return
	}

	hlog.Infof("SubmitVideo: userID=%d videoID=%d success", userID, req.VideoID)
	c.JSON(200, hertz.Success(map[string]interface{}{
		"video_id": resp.VideoId,
		"success":  true,
	}))
}

// @Summary 发布视频
// @Description 发布视频（需要认证），支持创建新视频或更新已有视频并发布
// @Tags video
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body PublishVideoRequest true "发布信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/publish [post]
func (h *VideoHandler) PublishVideo(ctx context.Context, c *app.RequestContext) {
	var req PublishVideoRequest
	if err := c.Bind(&req); err != nil {
		hlog.Errorf("PublishVideo: invalid request: %v", err)
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	resp, err := h.videoClient.PublishVideo(ctx, req.VideoID)
	if err != nil {
		hlog.Errorf("PublishVideo: videoID=%d failed: %v", req.VideoID, err)
		c.JSON(500, hertz.Fail(500, "Failed to publish video"))
		return
	}

	hlog.Infof("PublishVideo: videoID=%d success", req.VideoID)
	c.JSON(200, hertz.Success(map[string]interface{}{
		"video_id": resp.VideoId,
		"success":  true,
	}))
}

// @Summary 获取发布列表
// @Description 获取指定用户的发布视频列表，支持游标分页
// @Tags video
// @Produce json
// @Param user_id query int64 true "用户ID"
// @Param limit query int32 false "返回数量" default(10)
// @Param cursor query string false "游标分页（首次不传）" default("")
// @Success 200 {object} SuccessResponse{data=HotVideoListResponse}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/list/published [get]
func (h *VideoHandler) GetPublishedList(ctx context.Context, c *app.RequestContext) {
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil || userID == 0 {
		c.JSON(400, hertz.Fail(400, "Invalid user_id"))
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 32)
	cursor := c.Query("cursor")

	resp, err := h.videoClient.GetPublishedList(ctx, userID, int32(limit), cursor)
	if err != nil {
		hlog.Errorf("GetPublishedList: userID=%d failed: %v", userID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get published list"))
		return
	}

	videos := make([]*HotVideoItem, 0, len(resp.Videos))
	for _, v := range resp.Videos {
		videos = append(videos, &HotVideoItem{
			ID:           v.Id,
			Title:        v.Title,
			CoverURL:     v.CoverUrl,
			Duration:     int(v.Duration),
			ViewCount:    v.ViewCount,
			CommentCount: v.CommentCount,
		})
	}

	c.JSON(200, hertz.Success(HotVideoListResponse{
		Videos:     videos,
		NextCursor: resp.NextCursor,
		HasMore:    resp.HasMore,
	}))
}

// @Summary 获取上传凭证
// @Description 获取七牛云上传凭证，前端拿到 token 后直传到七牛云（需要认证），upload_type: 1=视频, 2=封面
// @Tags video
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body handler.UploadTokenRequest true "video_id 和 upload_type"
// @Success 200 {object} UploadTokenResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/upload/token [post]
func (h *VideoHandler) GetUploadToken(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64(constants.ContextKeyUserID)
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req UploadTokenRequest
	if err := c.Bind(&req); err != nil {
		hlog.Errorf("GetUploadToken: invalid request: %v", err)
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	if req.VideoID == 0 {
		c.JSON(400, hertz.Fail(400, "video_id is required"))
		return
	}
	if req.UploadType != constants.UploadTokenTypeVideo && req.UploadType != constants.UploadTokenTypeCover {
		c.JSON(400, hertz.Fail(400, "invalid upload_type"))
		return
	}

	resp, err := h.videoClient.GetUploadToken(ctx, req.VideoID, req.UploadType)
	if err != nil {
		hlog.Errorf("GetUploadToken: userID=%d, videoID=%d, uploadType=%d failed: %v", userID, req.VideoID, req.UploadType, err)
		c.JSON(500, hertz.Fail(500, "Failed to get upload token"))
		return
	}

	c.JSON(200, hertz.Success(UploadTokenResponse{
		Token:  resp.Token,
		Key:    resp.Key,
		Domain: resp.Domain,
		Host:   resp.Host,
	}))
}
