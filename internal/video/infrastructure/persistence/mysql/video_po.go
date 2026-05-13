package mysql

import (
	"time"

	"gorm.io/gorm"

	constants "vicomova/pkg/constants"
)

type VideoPO struct {
	ID           int64          `gorm:"primaryKey;autoIncrement"`
	UserID       int64          `gorm:"not null;index:idx_user"`
	Title        string         `gorm:"size:255;not null"`
	Description  string         `gorm:"type:text"`
	CoverURL     string         `gorm:"size:512"`
	VideoURL     string         `gorm:"size:512;not null"`
	CategoryID   int            `gorm:"index:idx_category"`
	ViewCount    int64          `gorm:"default:0"`
	LikeCount    int64          `gorm:"default:0"`
	CommentCount int64          `gorm:"default:0"`
	Duration     int            `gorm:"default:0"`
	Status       int8           `gorm:"default:0;index:idx_status"`
	HotScore     float64        `gorm:"index:idx_hot_score"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (VideoPO) TableName() string {
	return constants.TableVideos
}
