package mysql

import (
	"vicomova/internal/interaction/domain/entity"
)

func LikeToPO(l *entity.Like) *LikePO {
	return &LikePO{
		ID:         l.ID,
		UserID:     l.UserID,
		TargetType: l.TargetType,
		TargetID:   l.TargetID,
	}
}

func POToLike(po *LikePO) *entity.Like {
	return &entity.Like{
		ID:         po.ID,
		UserID:     po.UserID,
		TargetType: po.TargetType,
		TargetID:   po.TargetID,
		CreatedAt:  po.CreatedAt,
	}
}

func CommentToPO(c *entity.Comment) *CommentPO {
	return &CommentPO{
		ID:        c.ID,
		UserID:    c.UserID,
		VideoID:   c.VideoID,
		ParentID:  c.ParentID,
		Content:   c.Content,
		LikeCount: c.LikeCount,
	}
}

func POToComment(po *CommentPO) *entity.Comment {
	return &entity.Comment{
		ID:        po.ID,
		UserID:    po.UserID,
		VideoID:   po.VideoID,
		ParentID:  po.ParentID,
		Content:   po.Content,
		LikeCount: po.LikeCount,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}
}

func FavoriteToPO(f *entity.Favorite) *FavoritePO {
	return &FavoritePO{
		ID:      f.ID,
		UserID:  f.UserID,
		VideoID: f.VideoID,
	}
}

func POToFavorite(po *FavoritePO) *entity.Favorite {
	return &entity.Favorite{
		ID:        po.ID,
		UserID:    po.UserID,
		VideoID:   po.VideoID,
		CreatedAt: po.CreatedAt,
	}
}
