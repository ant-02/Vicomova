package command

import (
	"context"

	"vicomova/internal/interaction/domain/entity"
)

func (s *InteractionCommandService) LikeVideo(ctx context.Context, cmd *LikeCommand) error {
	existing, err := s.likeRepo.Get(ctx, cmd.UserID, entity.TargetTypeVideo, cmd.TargetID)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}

	like := &entity.Like{
		UserID:     cmd.UserID,
		TargetType: entity.TargetTypeVideo,
		TargetID:   cmd.TargetID,
	}
	return s.likeRepo.Create(ctx, like)
}

func (s *InteractionCommandService) UnlikeVideo(ctx context.Context, cmd *LikeCommand) error {
	return s.likeRepo.Delete(ctx, cmd.UserID, entity.TargetTypeVideo, cmd.TargetID)
}

func (s *InteractionCommandService) LikeComment(ctx context.Context, cmd *LikeCommand) error {
	existing, err := s.likeRepo.Get(ctx, cmd.UserID, entity.TargetTypeComment, cmd.TargetID)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}

	like := &entity.Like{
		UserID:     cmd.UserID,
		TargetType: entity.TargetTypeComment,
		TargetID:   cmd.TargetID,
	}
	if err := s.likeRepo.Create(ctx, like); err != nil {
		return err
	}
	return s.commentRepo.IncrementLike(ctx, cmd.TargetID)
}

func (s *InteractionCommandService) AddFavorite(ctx context.Context, cmd *FavoriteCommand) error {
	existing, err := s.favoriteRepo.Get(ctx, cmd.UserID, cmd.VideoID)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}

	fav := &entity.Favorite{
		UserID:  cmd.UserID,
		VideoID: cmd.VideoID,
	}
	return s.favoriteRepo.Create(ctx, fav)
}

func (s *InteractionCommandService) RemoveFavorite(ctx context.Context, cmd *FavoriteCommand) error {
	return s.favoriteRepo.Delete(ctx, cmd.UserID, cmd.VideoID)
}

func (s *InteractionCommandService) Comment(ctx context.Context, cmd *CommentCommand) (*CommentResult, error) {
	comment := &entity.Comment{
		UserID:   cmd.UserID,
		VideoID:  cmd.VideoID,
		ParentID: cmd.ParentID,
		Content:  cmd.Content,
	}
	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, err
	}
	return &CommentResult{CommentID: comment.ID}, nil
}

func (s *InteractionCommandService) DeleteComment(ctx context.Context, userID, commentID int64) error {
	comment, err := s.commentRepo.Get(ctx, commentID)
	if err != nil {
		return err
	}
	if comment == nil {
		return nil
	}
	if comment.UserID != userID {
		return nil
	}
	return s.commentRepo.Delete(ctx, commentID)
}
