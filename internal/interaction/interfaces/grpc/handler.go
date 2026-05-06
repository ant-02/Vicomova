package grpc

import (
	"context"

	"vicomova/internal/interaction/application/command"
	"vicomova/internal/interaction/application/query"
	"vicomova/internal/interaction/domain/entity"
	interaction "vicomova/third_party/kitex_gen/interaction"
)

type InteractionHandler struct {
	cmdSvc *command.InteractionCommandService
	qrySvc *query.InteractionQueryService
}

func NewInteractionHandler(cmdSvc *command.InteractionCommandService, qrySvc *query.InteractionQueryService) *InteractionHandler {
	return &InteractionHandler{
		cmdSvc: cmdSvc,
		qrySvc: qrySvc,
	}
}

func (h *InteractionHandler) LikeVideo(ctx context.Context, req *interaction.LikeVideoRequest) (*interaction.LikeVideoResponse, error) {
	cmd := &command.LikeCommand{
		UserID:     req.UserId,
		TargetType: "video",
		TargetID:   req.VideoId,
	}
	err := h.cmdSvc.LikeVideo(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &interaction.LikeVideoResponse{}, nil
}

func (h *InteractionHandler) UnlikeVideo(ctx context.Context, req *interaction.UnlikeVideoRequest) (*interaction.UnlikeVideoResponse, error) {
	cmd := &command.LikeCommand{
		UserID:     req.UserId,
		TargetType: "video",
		TargetID:   req.VideoId,
	}
	err := h.cmdSvc.UnlikeVideo(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &interaction.UnlikeVideoResponse{}, nil
}

func (h *InteractionHandler) ListLikes(ctx context.Context, req *interaction.ListLikesRequest) (*interaction.ListLikesResponse, error) {
	result, err := h.qrySvc.ListLikes(ctx, req.UserId, req.TargetType, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}

	likes := result.Items.([]*entity.Like)
	protoLikes := make([]*interaction.Like, len(likes))
	for i, l := range likes {
		protoLikes[i] = &interaction.Like{
			Id:         l.ID,
			UserId:     l.UserID,
			TargetType: l.TargetType,
			TargetId:   l.TargetID,
			CreatedAt:  l.CreatedAt.Unix(),
		}
	}

	return &interaction.ListLikesResponse{
		Likes: protoLikes,
		Total: result.Total,
	}, nil
}

func (h *InteractionHandler) AddFavorite(ctx context.Context, req *interaction.AddFavoriteRequest) (*interaction.AddFavoriteResponse, error) {
	cmd := &command.FavoriteCommand{
		UserID:  req.UserId,
		VideoID: req.VideoId,
	}
	err := h.cmdSvc.AddFavorite(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &interaction.AddFavoriteResponse{}, nil
}

func (h *InteractionHandler) RemoveFavorite(ctx context.Context, req *interaction.RemoveFavoriteRequest) (*interaction.RemoveFavoriteResponse, error) {
	cmd := &command.FavoriteCommand{
		UserID:  req.UserId,
		VideoID: req.VideoId,
	}
	err := h.cmdSvc.RemoveFavorite(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &interaction.RemoveFavoriteResponse{}, nil
}

func (h *InteractionHandler) ListFavorites(ctx context.Context, req *interaction.ListFavoritesRequest) (*interaction.ListFavoritesResponse, error) {
	result, err := h.qrySvc.ListFavorites(ctx, req.UserId, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}

	favs := result.Items.([]*entity.Favorite)
	protoFavs := make([]*interaction.Favorite, len(favs))
	for i, f := range favs {
		protoFavs[i] = &interaction.Favorite{
			Id:        f.ID,
			UserId:    f.UserID,
			VideoId:   f.VideoID,
			CreatedAt: f.CreatedAt.Unix(),
		}
	}

	return &interaction.ListFavoritesResponse{
		Favorites: protoFavs,
		Total:     result.Total,
	}, nil
}

func (h *InteractionHandler) Comment(ctx context.Context, req *interaction.CommentRequest) (*interaction.CommentResponse, error) {
	cmd := &command.CommentCommand{
		UserID:   req.UserId,
		VideoID:  req.VideoId,
		ParentID: req.ParentId,
		Content:  req.Content,
	}
	result, err := h.cmdSvc.Comment(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &interaction.CommentResponse{CommentId: result.CommentID}, nil
}

func (h *InteractionHandler) DeleteComment(ctx context.Context, req *interaction.DeleteCommentRequest) (*interaction.DeleteCommentResponse, error) {
	err := h.cmdSvc.DeleteComment(ctx, req.UserId, req.CommentId)
	if err != nil {
		return nil, err
	}
	return &interaction.DeleteCommentResponse{}, nil
}

func (h *InteractionHandler) ListComments(ctx context.Context, req *interaction.ListCommentsRequest) (*interaction.ListCommentsResponse, error) {
	result, err := h.qrySvc.ListComments(ctx, req.VideoId, req.ParentId, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}

	comments := result.Items.([]*entity.Comment)
	protoComments := make([]*interaction.Comment, len(comments))
	for i, c := range comments {
		protoComments[i] = &interaction.Comment{
			Id:        c.ID,
			UserId:    c.UserID,
			VideoId:   c.VideoID,
			ParentId:  c.ParentID,
			Content:   c.Content,
			LikeCount: c.LikeCount,
			CreatedAt: c.CreatedAt.Unix(),
		}
	}

	return &interaction.ListCommentsResponse{
		Comments: protoComments,
		Total:    result.Total,
	}, nil
}

func (h *InteractionHandler) LikeComment(ctx context.Context, req *interaction.LikeCommentRequest) (*interaction.LikeCommentResponse, error) {
	cmd := &command.LikeCommand{
		UserID:     req.UserId,
		TargetType: "comment",
		TargetID:   req.CommentId,
	}
	err := h.cmdSvc.LikeComment(ctx, cmd)
	if err != nil {
		return nil, err
	}
	return &interaction.LikeCommentResponse{}, nil
}
