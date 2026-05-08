package handler

import (
	"context"
	"strconv"

	interactionRpc "vicomova/internal/interaction/interfaces/grpc"
	hertz "vicomova/pkg/infrastructure/hertz"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type InteractionHandler struct {
	interactionClient *interactionRpc.InteractionClient
}

func NewInteractionHandler(interactionClient *interactionRpc.InteractionClient) *InteractionHandler {
	return &InteractionHandler{interactionClient: interactionClient}
}

// @Summary 点赞视频
// @Description 点赞视频（需要认证）
// @Tags interaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body LikeRequest true "点赞请求"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/like [post]
func (h *InteractionHandler) LikeVideo(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req LikeRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	_, err := h.interactionClient.LikeVideo(ctx, userID, req.VideoID)
	if err != nil {
		hlog.Errorf("LikeVideo: userID=%d videoID=%d failed: %v", userID, req.VideoID, err)
		c.JSON(500, hertz.Fail(500, "Failed to like video"))
		return
	}

	hlog.Infof("LikeVideo: userID=%d videoID=%d success", userID, req.VideoID)
	c.JSON(200, hertz.Success(SuccessResponse{Success: true}))
}

// @Summary 取消点赞视频
// @Description 取消点赞视频（需要认证）
// @Tags interaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UnlikeRequest true "取消点赞请求"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/like [delete]
func (h *InteractionHandler) UnlikeVideo(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req UnlikeRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	_, err := h.interactionClient.UnlikeVideo(ctx, userID, req.VideoID)
	if err != nil {
		hlog.Errorf("UnlikeVideo: userID=%d videoID=%d failed: %v", userID, req.VideoID, err)
		c.JSON(500, hertz.Fail(500, "Failed to unlike video"))
		return
	}

	hlog.Infof("UnlikeVideo: userID=%d videoID=%d success", userID, req.VideoID)
	c.JSON(200, hertz.Success(SuccessResponse{Success: true}))
}

// @Summary 获取点赞列表
// @Description 获取用户点赞列表（需要认证）
// @Tags interaction
// @Produce json
// @Security BearerAuth
// @Param target_type query string false "目标类型 (video/comment)" default(video)
// @Param page query int32 false "页码" default(1)
// @Param size query int32 false "每页数量" default(10)
// @Success 200 {object} LikeListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/like/list [get]
func (h *InteractionHandler) ListLikes(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	targetType := c.DefaultQuery("target_type", "video")
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	size, _ := strconv.ParseInt(c.DefaultQuery("size", "10"), 10, 32)

	resp, err := h.interactionClient.ListLikes(ctx, userID, targetType, int32(page), int32(size))
	if err != nil {
		hlog.Errorf("ListLikes: userID=%d failed: %v", userID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get like list"))
		return
	}

	likes := make([]*LikeItem, 0, len(resp.Likes))
	for _, l := range resp.Likes {
		likes = append(likes, &LikeItem{
			ID:         l.Id,
			UserID:     l.UserId,
			TargetType: l.TargetType,
			TargetID:   l.TargetId,
			CreatedAt:  l.CreatedAt,
		})
	}

	c.JSON(200, hertz.Success(LikeListResponse{
		Likes: likes,
		Total: resp.Total,
	}))
}

// @Summary 收藏视频
// @Description 添加视频到收藏夹（需要认证）
// @Tags interaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body FavoriteRequest true "收藏请求"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/favorite [post]
func (h *InteractionHandler) AddFavorite(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req FavoriteRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	_, err := h.interactionClient.AddFavorite(ctx, userID, req.VideoID)
	if err != nil {
		hlog.Errorf("AddFavorite: userID=%d videoID=%d failed: %v", userID, req.VideoID, err)
		c.JSON(500, hertz.Fail(500, "Failed to add favorite"))
		return
	}

	hlog.Infof("AddFavorite: userID=%d videoID=%d success", userID, req.VideoID)
	c.JSON(200, hertz.Success(SuccessResponse{Success: true}))
}

// @Summary 取消收藏视频
// @Description 从收藏夹移除视频（需要认证）
// @Tags interaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UnfavoriteRequest true "取消收藏请求"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/favorite [delete]
func (h *InteractionHandler) RemoveFavorite(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req UnfavoriteRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	_, err := h.interactionClient.RemoveFavorite(ctx, userID, req.VideoID)
	if err != nil {
		hlog.Errorf("RemoveFavorite: userID=%d videoID=%d failed: %v", userID, req.VideoID, err)
		c.JSON(500, hertz.Fail(500, "Failed to remove favorite"))
		return
	}

	hlog.Infof("RemoveFavorite: userID=%d videoID=%d success", userID, req.VideoID)
	c.JSON(200, hertz.Success(SuccessResponse{Success: true}))
}

// @Summary 获取收藏列表
// @Description 获取用户收藏的视频列表（需要认证）
// @Tags interaction
// @Produce json
// @Security BearerAuth
// @Param page query int32 false "页码" default(1)
// @Param size query int32 false "每页数量" default(10)
// @Success 200 {object} FavoriteListResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/favorite/list [get]
func (h *InteractionHandler) ListFavorites(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	size, _ := strconv.ParseInt(c.DefaultQuery("size", "10"), 10, 32)

	resp, err := h.interactionClient.ListFavorites(ctx, userID, int32(page), int32(size))
	if err != nil {
		hlog.Errorf("ListFavorites: userID=%d failed: %v", userID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get favorite list"))
		return
	}

	favorites := make([]*FavoriteItem, 0, len(resp.Favorites))
	for _, f := range resp.Favorites {
		favorites = append(favorites, &FavoriteItem{
			ID:        f.Id,
			UserID:    f.UserId,
			VideoID:   f.VideoId,
			CreatedAt: f.CreatedAt,
		})
	}

	c.JSON(200, hertz.Success(FavoriteListResponse{
		Favorites: favorites,
		Total:     resp.Total,
	}))
}

// @Summary 评论视频
// @Description 对视频进行评论（需要认证）
// @Tags interaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CommentRequest true "评论请求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/comment [post]
func (h *InteractionHandler) Comment(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req CommentRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	resp, err := h.interactionClient.Comment(ctx, userID, req.VideoID, req.ParentID, req.Content)
	if err != nil {
		hlog.Errorf("Comment: userID=%d videoID=%d failed: %v", userID, req.VideoID, err)
		c.JSON(500, hertz.Fail(500, "Failed to comment"))
		return
	}

	hlog.Infof("Comment: userID=%d videoID=%d success, commentID=%d", userID, req.VideoID, resp.CommentId)
	c.JSON(200, hertz.Success(map[string]interface{}{
		"comment_id": resp.CommentId,
		"success":    true,
	}))
}

// @Summary 删除评论
// @Description 删除自己的评论（需要认证）
// @Tags interaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body DeleteCommentRequest true "删除评论请求"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/comment [delete]
func (h *InteractionHandler) DeleteComment(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req DeleteCommentRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	_, err := h.interactionClient.DeleteComment(ctx, userID, req.CommentID)
	if err != nil {
		hlog.Errorf("DeleteComment: userID=%d commentID=%d failed: %v", userID, req.CommentID, err)
		c.JSON(500, hertz.Fail(500, "Failed to delete comment"))
		return
	}

	hlog.Infof("DeleteComment: userID=%d commentID=%d success", userID, req.CommentID)
	c.JSON(200, hertz.Success(SuccessResponse{Success: true}))
}

// @Summary 获取评论列表
// @Description 获取视频的评论列表
// @Tags interaction
// @Produce json
// @Param video_id query int64 true "视频ID"
// @Param parent_id query int64 false "父评论ID (0表示根评论)" default(0)
// @Param page query int32 false "页码" default(1)
// @Param size query int32 false "每页数量" default(10)
// @Success 200 {object} CommentListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /video/comment/list [get]
func (h *InteractionHandler) ListComments(ctx context.Context, c *app.RequestContext) {
	videoID, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid video_id"))
		return
	}

	parentID, _ := strconv.ParseInt(c.DefaultQuery("parent_id", "0"), 10, 64)
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 32)
	size, _ := strconv.ParseInt(c.DefaultQuery("size", "10"), 10, 32)

	resp, err := h.interactionClient.ListComments(ctx, videoID, parentID, int32(page), int32(size))
	if err != nil {
		hlog.Errorf("ListComments: videoID=%d failed: %v", videoID, err)
		c.JSON(500, hertz.Fail(500, "Failed to get comment list"))
		return
	}

	comments := make([]*CommentItem, 0, len(resp.Comments))
	for _, c := range resp.Comments {
		comments = append(comments, &CommentItem{
			ID:        c.Id,
			UserID:    c.UserId,
			VideoID:   c.VideoId,
			ParentID:  c.ParentId,
			Content:   c.Content,
			LikeCount: c.LikeCount,
			CreatedAt: c.CreatedAt,
		})
	}

	c.JSON(200, hertz.Success(CommentListResponse{
		Comments: comments,
		Total:    resp.Total,
	}))
}

// @Summary 点赞评论
// @Description 点赞评论（需要认证）
// @Tags interaction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body LikeCommentRequest true "点赞评论请求"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /comment/like [post]
func (h *InteractionHandler) LikeComment(ctx context.Context, c *app.RequestContext) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		c.JSON(401, hertz.Fail(401, "Unauthorized"))
		return
	}

	var req LikeCommentRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(400, hertz.Fail(400, "Invalid request body"))
		return
	}

	_, err := h.interactionClient.LikeComment(ctx, userID, req.CommentID)
	if err != nil {
		hlog.Errorf("LikeComment: userID=%d commentID=%d failed: %v", userID, req.CommentID, err)
		c.JSON(500, hertz.Fail(500, "Failed to like comment"))
		return
	}

	hlog.Infof("LikeComment: userID=%d commentID=%d success", userID, req.CommentID)
	c.JSON(200, hertz.Success(SuccessResponse{Success: true}))
}
