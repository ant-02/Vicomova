package grpc

import (
	"context"

	"vicomova/internal/video/application/command"
	"vicomova/internal/video/application/query"
	"vicomova/internal/video/domain/repository"
	video "vicomova/third_party/kitex_gen/video"
)

type VideoHandler struct {
	cmdSvc            *command.VideoCommandService
	qrySvc            *query.VideoQueryService
	viewCountProducer repository.ViewCountProducer
}

func NewVideoHandler(cmdSvc *command.VideoCommandService, qrySvc *query.VideoQueryService, viewCountProducer repository.ViewCountProducer) *VideoHandler {
	return &VideoHandler{
		cmdSvc:            cmdSvc,
		qrySvc:            qrySvc,
		viewCountProducer: viewCountProducer,
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

func (h *VideoHandler) ListCategoryVideos(ctx context.Context, req *video.ListCategoryVideosRequest) (*video.ListCategoryVideosResponse, error) {
	result, err := h.qrySvc.ListCategoryVideos(ctx, &query.CategoryVideosQuery{
		CategoryID: int(req.CategoryId),
		Cursor:     req.Cursor,
		Limit:      int(req.Limit),
	})
	if err != nil {
		return nil, err
	}

	protoVideos := make([]*video.Video, len(result.Videos))
	for i, item := range result.Videos {
		protoVideos[i] = &video.Video{
			Id:           item.VideoID,
			Title:        item.Title,
			CoverUrl:     item.CoverURL,
			Duration:     int32(item.Duration),
			ViewCount:    item.ViewCount,
			CommentCount: item.CommentCount,
		}
	}

	return &video.ListCategoryVideosResponse{
		Videos:     protoVideos,
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}

func (h *VideoHandler) ListHotVideos(ctx context.Context, req *video.ListHotVideosRequest) (*video.ListHotVideosResponse, error) {
	result, err := h.qrySvc.ListHotVideos(ctx, &query.ListHotVideosQuery{
		Cursor: req.Cursor,
		Limit:  int(req.Limit),
	})
	if err != nil {
		return nil, err
	}

	protoVideos := make([]*video.Video, len(result.Videos))
	for i, item := range result.Videos {
		protoVideos[i] = &video.Video{
			Id:           item.VideoID,
			Title:        item.Title,
			CoverUrl:     item.CoverURL,
			Duration:     int32(item.Duration),
			ViewCount:    item.ViewCount,
			CommentCount: item.CommentCount,
			UserName:     item.UserName,
		}
	}

	return &video.ListHotVideosResponse{
		Videos:     protoVideos,
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}

func (h *VideoHandler) GetPublishedList(ctx context.Context, req *video.GetPublishedListRequest) (*video.GetPublishedListResponse, error) {
	result, err := h.qrySvc.ListPublishedVideos(ctx, req.UserId, &query.PublishedVideosQuery{
		Cursor: req.Cursor,
		Limit:  int(req.Limit),
	})
	if err != nil {
		return nil, err
	}

	protoVideos := make([]*video.Video, len(result.Videos))
	for i, item := range result.Videos {
		protoVideos[i] = &video.Video{
			Id:           item.VideoID,
			Title:        item.Title,
			CoverUrl:     item.CoverURL,
			Duration:     int32(item.Duration),
			ViewCount:    item.ViewCount,
			CommentCount: item.CommentCount,
		}
	}

	return &video.GetPublishedListResponse{
		Videos:     protoVideos,
		NextCursor: result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}

func (h *VideoHandler) IncrementView(ctx context.Context, req *video.IncrementViewRequest) (*video.IncrementViewResponse, error) {
	if h.viewCountProducer != nil {
		_ = h.viewCountProducer.Record(ctx, req.VideoId)
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
