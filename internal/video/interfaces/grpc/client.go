package grpc

import (
	"context"

	video "vicomova/third_party/kitex_gen/video"
	videoservice "vicomova/third_party/kitex_gen/video/videoservice"

	"github.com/cloudwego/kitex/client"
)

type VideoClient struct {
	cli videoservice.Client
}

// NewVideoClient creates RPC client, directly connects to specified address
func NewVideoClient(serviceName, addr string) (*VideoClient, error) {
	cli, err := videoservice.NewClient(serviceName,
		client.WithHostPorts(addr),
	)
	if err != nil {
		return nil, err
	}
	return &VideoClient{cli: cli}, nil
}

func (c *VideoClient) SaveVideo(ctx context.Context, videoID int64, userID int64, title, description, coverURL, videoURL string, duration int32, categoryID int32) (*video.SaveVideoResponse, error) {
	return c.cli.SaveVideo(ctx, &video.SaveVideoRequest{
		VideoId:     videoID,
		UserId:      userID,
		Title:       title,
		Description: description,
		CategoryId:  categoryID,
		CoverUrl:    coverURL,
		VideoUrl:    videoURL,
		Duration:    duration,
	})
}

func (c *VideoClient) SubmitVideo(ctx context.Context, videoID int64, userID int64) (*video.SubmitVideoResponse, error) {
	return c.cli.SubmitVideo(ctx, &video.SubmitVideoRequest{
		VideoId: videoID,
		UserId:  userID,
	})
}

func (c *VideoClient) PublishVideo(ctx context.Context, videoID int64) (*video.PublishVideoResponse, error) {
	return c.cli.PublishVideo(ctx, &video.PublishVideoRequest{
		VideoId: videoID,
	})
}

func (c *VideoClient) GetVideoStream(ctx context.Context, videoID int64) (*video.GetVideoStreamResponse, error) {
	return c.cli.GetVideoStream(ctx, &video.GetVideoStreamRequest{
		VideoId: videoID,
	})
}

func (c *VideoClient) ListCategoryVideos(ctx context.Context, categoryID int32, limit int32, cursor string) (*video.ListCategoryVideosResponse, error) {
	return c.cli.ListCategoryVideos(ctx, &video.ListCategoryVideosRequest{
		CategoryId: categoryID,
		Limit:      limit,
		Cursor:     cursor,
	})
}

func (c *VideoClient) ListHotVideos(ctx context.Context, limit int32, cursor string) (*video.ListHotVideosResponse, error) {
	return c.cli.ListHotVideos(ctx, &video.ListHotVideosRequest{
		Limit:  limit,
		Cursor: cursor,
	})
}

func (c *VideoClient) GetPublishedList(ctx context.Context, userID int64, limit int32, cursor string) (*video.GetPublishedListResponse, error) {
	return c.cli.GetPublishedList(ctx, &video.GetPublishedListRequest{
		UserId: userID,
		Limit:  limit,
		Cursor: cursor,
	})
}

func (c *VideoClient) IncrementView(ctx context.Context, videoID int64) (*video.IncrementViewResponse, error) {
	return c.cli.IncrementView(ctx, &video.IncrementViewRequest{
		VideoId: videoID,
	})
}

func (c *VideoClient) GetUploadToken(ctx context.Context, videoID int64, uploadType int32) (*video.GetUploadTokenResponse, error) {
	return c.cli.GetUploadToken(ctx, &video.GetUploadTokenRequest{
		VideoId:    videoID,
		UploadType: uploadType,
	})
}
