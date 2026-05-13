package mysql

import (
	"time"

	"gorm.io/gorm"
)

type CategoryPO struct {
	ID        int            `gorm:"primaryKey;autoIncrement"`
	Name      string         `gorm:"size:64;not null"`
	ParentID  int            `gorm:"default:0"`
	SortOrder int            `gorm:"default:0"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (CategoryPO) TableName() string {
	return "categories"
}
