package repository

import (
    "context"

    "vicomova/internal/interaction/domain/entity"
)

// LikeRepository 点赞仓储
type LikeRepository interface {
    Create(ctx context.Context, like *entity.Like) error
    Delete(ctx context.Context, userID int64, targetType string, targetID int64) error
    Get(ctx context.Context, userID int64, targetType string, targetID int64) (*entity.Like, error)
    ListByUser(ctx context.Context, userID int64, targetType string, page, size int) ([]*entity.Like, int64, error)
    Count(ctx context.Context, targetType string, targetID int64) (int64, error)
}

// CommentRepository 评论仓储
type CommentRepository interface {
    Create(ctx context.Context, comment *entity.Comment) error
    Delete(ctx context.Context, id int64) error
    Get(ctx context.Context, id int64) (*entity.Comment, error)
    ListByVideo(ctx context.Context, videoID int64, page, size int) ([]*entity.Comment, int64, error)
    ListByParent(ctx context.Context, parentID int64, page, size int) ([]*entity.Comment, int64, error)
    IncrementLike(ctx context.Context, id int64) error
}

// FavoriteRepository 收藏仓储
type FavoriteRepository interface {
    Create(ctx context.Context, fav *entity.Favorite) error
    Delete(ctx context.Context, userID, videoID int64) error
    Get(ctx context.Context, userID, videoID int64) (*entity.Favorite, error)
    ListByUser(ctx context.Context, userID int64, page, size int) ([]*entity.Favorite, int64, error)
}
