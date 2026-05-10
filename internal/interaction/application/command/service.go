package command

import (
	"vicomova/internal/interaction/domain/repository"
)

type InteractionCommandService struct {
	likeRepo     repository.LikeRepository
	commentRepo  repository.CommentRepository
	favoriteRepo repository.FavoriteRepository
}

func NewInteractionCommandService(
	likeRepo repository.LikeRepository,
	commentRepo repository.CommentRepository,
	favoriteRepo repository.FavoriteRepository,
) *InteractionCommandService {
	return &InteractionCommandService{
		likeRepo:     likeRepo,
		commentRepo:  commentRepo,
		favoriteRepo: favoriteRepo,
	}
}
