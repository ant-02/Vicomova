package mysql

import (
	"time"
)

type OfflineMessagePO struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	ToUserID   int64     `gorm:"not null;index:idx_to_user"`
	FromUserID int64     `gorm:"not null;index:idx_from_user"`
	Content    string    `gorm:"type:text;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	IsRead     bool      `gorm:"default:false;index:idx_is_read"`
}

func (OfflineMessagePO) TableName() string {
	return "offline_messages"
}
