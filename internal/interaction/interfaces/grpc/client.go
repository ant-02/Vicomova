package grpc

import (
	"context"

	interaction "vicomova/third_party/kitex_gen/interaction"
	interactionservice "vicomova/third_party/kitex_gen/interaction/interactionservice"

	"github.com/cloudwego/kitex/client"
)

type InteractionClient struct {
	cli interactionservice.Client
}

// NewInteractionClient creates RPC client, directly connects to specified address
func NewInteractionClient(serviceName, addr string) (*InteractionClient, error) {
	cli, err := interactionservice.NewClient(serviceName,
		client.WithHostPorts(addr),
	)
	if err != nil {
		return nil, err
	}
	return &InteractionClient{cli: cli}, nil
}

func (c *InteractionClient) LikeVideo(ctx context.Context, userID, videoID int64) (*interaction.LikeVideoResponse, error) {
	return c.cli.LikeVideo(ctx, &interaction.LikeVideoRequest{
		UserId:  userID,
		VideoId: videoID,
	})
}

func (c *InteractionClient) UnlikeVideo(ctx context.Context, userID, videoID int64) (*interaction.UnlikeVideoResponse, error) {
	return c.cli.UnlikeVideo(ctx, &interaction.UnlikeVideoRequest{
		UserId:  userID,
		VideoId: videoID,
	})
}

func (c *InteractionClient) ListLikes(ctx context.Context, userID int64, targetType string, cursor int64, limit int32) (*interaction.ListLikesResponse, error) {
	return c.cli.ListLikes(ctx, &interaction.ListLikesRequest{
		UserId:     userID,
		TargetType: targetType,
		Cursor:     cursor,
		Limit:      limit,
	})
}

func (c *InteractionClient) AddFavorite(ctx context.Context, userID, videoID int64) (*interaction.AddFavoriteResponse, error) {
	return c.cli.AddFavorite(ctx, &interaction.AddFavoriteRequest{
		UserId:  userID,
		VideoId: videoID,
	})
}

func (c *InteractionClient) RemoveFavorite(ctx context.Context, userID, videoID int64) (*interaction.RemoveFavoriteResponse, error) {
	return c.cli.RemoveFavorite(ctx, &interaction.RemoveFavoriteRequest{
		UserId:  userID,
		VideoId: videoID,
	})
}

func (c *InteractionClient) ListFavorites(ctx context.Context, userID int64, cursor int64, limit int32) (*interaction.ListFavoritesResponse, error) {
	return c.cli.ListFavorites(ctx, &interaction.ListFavoritesRequest{
		UserId: userID,
		Cursor: cursor,
		Limit:  limit,
	})
}

func (c *InteractionClient) Comment(ctx context.Context, userID, videoID, parentID int64, content string) (*interaction.CommentResponse, error) {
	return c.cli.Comment(ctx, &interaction.CommentRequest{
		UserId:   userID,
		VideoId:  videoID,
		ParentId: parentID,
		Content:  content,
	})
}

func (c *InteractionClient) DeleteComment(ctx context.Context, userID, commentID int64) (*interaction.DeleteCommentResponse, error) {
	return c.cli.DeleteComment(ctx, &interaction.DeleteCommentRequest{
		UserId:    userID,
		CommentId: commentID,
	})
}

func (c *InteractionClient) ListComments(ctx context.Context, videoID int64, parentID int64, cursor int64, limit int32) (*interaction.ListCommentsResponse, error) {
	return c.cli.ListComments(ctx, &interaction.ListCommentsRequest{
		VideoId:  videoID,
		ParentId: parentID,
		Cursor:   cursor,
		Limit:    limit,
	})
}

func (c *InteractionClient) LikeComment(ctx context.Context, userID, commentID int64) (*interaction.LikeCommentResponse, error) {
	return c.cli.LikeComment(ctx, &interaction.LikeCommentRequest{
		UserId:    userID,
		CommentId: commentID,
	})
}
