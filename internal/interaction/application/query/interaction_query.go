package query

import (
    "context"

    "vicomova/internal/interaction/domain/entity"
    "vicomova/internal/interaction/domain/repository"
)

type ListResult struct {
    Items interface{}
    Total int64
}

type InteractionQueryService struct {
    likeRepo    repository.LikeRepository
    commentRepo repository.CommentRepository
    favoriteRepo repository.FavoriteRepository
}

func NewInteractionQueryService(
    likeRepo repository.LikeRepository,
    commentRepo repository.CommentRepository,
    favoriteRepo repository.FavoriteRepository,
) *InteractionQueryService {
    return &InteractionQueryService{
        likeRepo:    likeRepo,
        commentRepo: commentRepo,
        favoriteRepo: favoriteRepo,
    }
}

func (s *InteractionQueryService) ListLikes(ctx context.Context, userID int64, targetType string, page, size int) (*ListResult, error) {
    likes, total, err := s.likeRepo.ListByUser(ctx, userID, targetType, page, size)
    if err != nil {
        return nil, err
    }
    return &ListResult{Items: likes, Total: total}, nil
}

func (s *InteractionQueryService) ListFavorites(ctx context.Context, userID int64, page, size int) (*ListResult, error) {
    favs, total, err := s.favoriteRepo.ListByUser(ctx, userID, page, size)
    if err != nil {
        return nil, err
    }
    return &ListResult{Items: favs, Total: total}, nil
}

func (s *InteractionQueryService) ListComments(ctx context.Context, videoID int64, parentID int64, page, size int) (*ListResult, error) {
    var comments []*entity.Comment
    var total int64
    var err error

    if parentID == 0 {
        comments, total, err = s.commentRepo.ListByVideo(ctx, videoID, page, size)
    } else {
        comments, total, err = s.commentRepo.ListByParent(ctx, parentID, page, size)
    }
    if err != nil {
        return nil, err
    }
    return &ListResult{Items: comments, Total: total}, nil
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
