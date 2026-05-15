package entity

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID             int64          `json:"id"`
	ConversationID int64          `json:"conversation_id"`
	SenderID       int64          `json:"sender_id"`
	Content        string         `json:"content"`
	CreatedAt      time.Time      `json:"created_at"`
	DeletedAt      gorm.DeletedAt `json:"-"`
}
