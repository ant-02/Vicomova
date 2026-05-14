package query

import (
	"context"

	"vicomova/internal/interaction/domain/entity"
)

func (s *InteractionQueryService) ListLikes(ctx context.Context, userID int64, targetType string, cursor int64, limit int) (*ListResult, bool, error) {
	likes, hasMore, err := s.likeRepo.ListByUser(ctx, userID, targetType, cursor, limit)
	if err != nil {
		return nil, false, err
	}
	return &ListResult{Items: likes}, hasMore, nil
}

func (s *InteractionQueryService) ListFavorites(ctx context.Context, userID int64, cursor int64, limit int) (*ListResult, bool, error) {
	favs, hasMore, err := s.favoriteRepo.ListByUser(ctx, userID, cursor, limit)
	if err != nil {
		return nil, false, err
	}
	return &ListResult{Items: favs}, hasMore, nil
}

func (s *InteractionQueryService) ListComments(ctx context.Context, videoID int64, parentID int64, cursor int64, limit int) (*ListResult, bool, error) {
	var comments []*entity.Comment
	var hasMore bool
	var err error

	if parentID == 0 {
		comments, hasMore, err = s.commentRepo.ListByVideo(ctx, videoID, cursor, limit)
	} else {
		comments, hasMore, err = s.commentRepo.ListByParent(ctx, parentID, cursor, limit)
	}
	if err != nil {
		return nil, false, err
	}
	return &ListResult{Items: comments}, hasMore, nil
}

func (s *InteractionQueryService) IsLiked(ctx context.Context, userID int64, targetType string, targetID int64) (bool, error) {
	like, err := s.likeRepo.Get(ctx, userID, targetType, targetID)
	if err != nil {
		return false, err
	}
	return like != nil, nil
}

func (s *InteractionQueryService) IsFavorited(ctx context.Context, userID, videoID int64) (bool, error) {
	fav, err := s.favoriteRepo.Get(ctx, userID, videoID)
	if err != nil {
		return false, err
	}
	return fav != nil, nil
}
