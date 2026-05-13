package entity

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	ParentID  int            `json:"parent_id"`
	SortOrder int            `json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
