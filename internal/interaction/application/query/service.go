package query

import (
	"vicomova/internal/interaction/domain/repository"
)

type InteractionQueryService struct {
	likeRepo     repository.LikeRepository
	commentRepo  repository.CommentRepository
	favoriteRepo repository.FavoriteRepository
}

func NewInteractionQueryService(
	likeRepo repository.LikeRepository,
	commentRepo repository.CommentRepository,
	favoriteRepo repository.FavoriteRepository,
) *InteractionQueryService {
	return &InteractionQueryService{
		likeRepo:     likeRepo,
		commentRepo:  commentRepo,
		favoriteRepo: favoriteRepo,
	}
}
