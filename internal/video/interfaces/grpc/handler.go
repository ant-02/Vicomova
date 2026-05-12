package grpc

import (
	"context"

	"vicomova/internal/video/application/command"
	"vicomova/internal/video/application/query"
	"vicomova/internal/video/domain/entity"
	video "vicomova/third_party/kitex_gen/video"
)

type VideoHandler struct {
	cmdSvc *command.VideoCommandService
	qrySvc *query.VideoQueryService
}

func NewVideoHandler(cmdSvc *command.VideoCommandService, qrySvc *query.VideoQueryService) *VideoHandler {
	return &VideoHandler{
		cmdSvc: cmdSvc,
		qrySvc: qrySvc,
	}
}

func (h *VideoHandler) SaveVideo(ctx context.Context, req *video.SaveVideoRequest) (*video.SaveVideoResponse, error) {
	cmd := &command.SaveVideoCommand{
		VideoID:     req.VideoId,
		UserID:      req.UserId,
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  int(req.CategoryId),
		CoverURL:    req.CoverUrl,
		VideoURL:    req.VideoUrl,
		Duration:    int(req.Duration),
	}

	result, err := h.cmdSvc.Save(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &video.SaveVideoResponse{VideoId: result.VideoID}, nil
}

func (h *VideoHandler) SubmitVideo(ctx context.Context, req *video.SubmitVideoRequest) (*video.SubmitVideoResponse, error) {
	cmd := &command.SubmitVideoCommand{
		VideoID: req.VideoId,
		UserID:  req.UserId,
	}

	result, err := h.cmdSvc.Submit(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &video.SubmitVideoResponse{VideoId: result.VideoID}, nil
}

func (h *VideoHandler) PublishVideo(ctx context.Context, req *video.PublishVideoRequest) (*video.PublishVideoResponse, error) {
	cmd := &command.PublishVideoCommand{
		VideoID: req.VideoId,
	}

	result, err := h.cmdSvc.Publish(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &video.PublishVideoResponse{VideoId: result.VideoID}, nil
}

func (h *VideoHandler) GetVideoStream(ctx context.Context, req *video.GetVideoStreamRequest) (*video.GetVideoStreamResponse, error) {
	v, err := h.qrySvc.GetVideoStream(ctx, req.VideoId)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}

	return &video.GetVideoStreamResponse{
		Id:           v.Video.ID,
		UserId:       v.Video.UserID,
		Title:        v.Video.Title,
		Description:  v.Video.Description,
		CoverUrl:     v.Video.CoverURL,
		VideoUrl:     v.Video.VideoURL,
		CategoryId:   int32(v.Video.CategoryID),
		ViewCount:    v.Video.ViewCount,
		LikeCount:    v.Video.LikeCount,
		CommentCount: v.Video.CommentCount,
		Duration:     int32(v.Video.Duration),
	}, nil
}

func (h *VideoHandler) ListByCategory(ctx context.Context, req *video.ListByCategoryRequest) (*video.ListByCategoryResponse, error) {
	result, err := h.qrySvc.ListByCategory(ctx, int(req.CategoryId), int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}

	videos := make([]*video.Video, len(result.Videos))
	for i, v := range result.Videos {
		videos[i] = toProtoVideo(v)
	}

	return &video.ListByCategoryResponse{
		Videos: videos,
		Total:  result.Total,
	}, nil
}

func (h *VideoHandler) ListHotVideos(ctx context.Context, req *video.ListHotVideosRequest) (*video.ListHotVideosResponse, error) {
	videos, err := h.qrySvc.ListHot(ctx, int(req.Limit))
	if err != nil {
		return nil, err
	}

	protoVideos := make([]*video.Video, len(videos))
	for i, v := range videos {
		protoVideos[i] = toProtoVideo(v)
	}

	return &video.ListHotVideosResponse{Videos: protoVideos}, nil
}

func (h *VideoHandler) GetPublishedList(ctx context.Context, req *video.GetPublishedListRequest) (*video.GetPublishedListResponse, error) {
	result, err := h.qrySvc.ListByUser(ctx, req.UserId, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}

	videos := make([]*video.Video, len(result.Videos))
	for i, v := range result.Videos {
		videos[i] = toProtoVideo(v)
	}

	return &video.GetPublishedListResponse{
		Videos: videos,
		Total:  result.Total,
	}, nil
}

func (h *VideoHandler) IncrementView(ctx context.Context, req *video.IncrementViewRequest) (*video.IncrementViewResponse, error) {
	err := h.qrySvc.IncrementView(ctx, req.VideoId)
	if err != nil {
		return nil, err
	}
	return &video.IncrementViewResponse{}, nil
}

func (h *VideoHandler) GetUploadToken(ctx context.Context, req *video.GetUploadTokenRequest) (*video.GetUploadTokenResponse, error) {
	result, err := h.qrySvc.GetUploadToken(ctx, req.VideoId, req.UploadType)
	if err != nil {
		return nil, err
	}

	return &video.GetUploadTokenResponse{
		Token:  result.Token,
		Key:    result.Key,
		Domain: result.Domain,
		Host:   result.Host,
	}, nil
}

func toProtoVideo(v *entity.Video) *video.Video {
	return &video.Video{
		Id:           v.ID,
		UserId:       v.UserID,
		Title:        v.Title,
		Description:  v.Description,
		CoverUrl:     v.CoverURL,
		VideoUrl:     v.VideoURL,
		CategoryId:   int32(v.CategoryID),
		ViewCount:    v.ViewCount,
		LikeCount:    v.LikeCount,
		CommentCount: v.CommentCount,
		Duration:     int32(v.Duration),
		CreatedAt:    v.CreatedAt.Unix(),
	}
}
