package mysql

import (
	"context"

	constants "vicomova/pkg/constants"
)

type LikePO struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	UserID     int64     `gorm:"not null;uniqueIndex:uk_like"`
	TargetType string    `gorm:"size:32;not null;uniqueIndex:uk_like"`
	TargetID   int64     `gorm:"not null;uniqueIndex:uk_like;index:idx_target"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func (LikePO) TableName() string {
	return constants.TableLikes
}

type CommentPO struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserID    int64     `gorm:"not null;index:idx_user"`
	VideoID   int64     `gorm:"not null;index:idx_video"`
	ParentID  int64     `gorm:"default:0;index:idx_parent"`
	Content   string    `gorm:"type:text;not null"`
	LikeCount int64     `gorm:"default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (CommentPO) TableName() string {
	return constants.TableComments
}

type FavoritePO struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserID    int64     `gorm:"not null;uniqueIndex:uk_favorite"`
	VideoID   int64     `gorm:"not null;uniqueIndex:uk_favorite;index:idx_video"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (FavoritePO) TableName() string {
	return constants.TableFavorites
}
