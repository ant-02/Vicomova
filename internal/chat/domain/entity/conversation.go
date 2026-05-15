package entity

import (
	"time"

	"gorm.io/gorm"
)

const (
	ConversationType1v1   = 1
	ConversationTypeGroup = 2
)

type Conversation struct {
	ID        int64          `json:"id"`
	Type      int8           `json:"type"` // 1=1v1, 2=group
	Name      *string        `json:"name,omitempty"`
	Avatar    *string        `json:"avatar,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
