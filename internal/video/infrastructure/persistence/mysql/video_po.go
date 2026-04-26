package mysql

import "time"

type VideoPO struct {
    ID           int64     `gorm:"primaryKey;autoIncrement"`
    UserID       int64     `gorm:"not null;index:idx_user"`
    Title        string    `gorm:"size:255;not null"`
    Description  string    `gorm:"type:text"`
    CoverURL     string    `gorm:"size:512"`
    VideoURL     string    `gorm:"size:512;not null"`
    CategoryID   int       `gorm:"index:idx_category"`
    ViewCount    int64     `gorm:"default:0"`
    LikeCount    int64     `gorm:"default:0"`
    CommentCount int64     `gorm:"default:0"`
    Duration     int       `gorm:"default:0"`
    Status       int8      `gorm:"default:0;index:idx_status"`
    CreatedAt    time.Time `gorm:"autoCreateTime"`
    UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (VideoPO) TableName() string {
    return "videos"
}
