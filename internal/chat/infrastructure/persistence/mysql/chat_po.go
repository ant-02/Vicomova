package mysql

import (
	"time"

	"gorm.io/gorm"
)

type FollowPO struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	FollowerID  int64          `gorm:"not null;uniqueIndex:uk_follow;index:idx_following"`
	FollowingID int64          `gorm:"not null;uniqueIndex:uk_follow"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (FollowPO) TableName() string {
	return "user_follows"
}

type ConversationPO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	Type      int8           `gorm:"not null;index:idx_type"`
	Name      *string        `gorm:"size:128"`
	Avatar    *string        `gorm:"size:512"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (ConversationPO) TableName() string {
	return "conversations"
}

type GroupMemberPO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	GroupID   int64          `gorm:"not null;uniqueIndex:uk_member;index:idx_user"`
	UserID    int64          `gorm:"not null;uniqueIndex:uk_member"`
	Role      int8           `gorm:"default:0"`
	JoinedAt  time.Time      `gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (GroupMemberPO) TableName() string {
	return "group_members"
}

type MessagePO struct {
	ID             int64          `gorm:"primaryKey;autoIncrement"`
	ConversationID int64          `gorm:"not null;index:idx_conversation"`
	SenderID       int64          `gorm:"not null;index:idx_sender"`
	Content        string         `gorm:"type:text;not null"`
	CreatedAt      time.Time      `gorm:"autoCreateTime;index:idx_created"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (MessagePO) TableName() string {
	return "messages"
}
