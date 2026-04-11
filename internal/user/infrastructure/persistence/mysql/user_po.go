package mysql

import "time"

type UserPO struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"uniqueIndex;size:64;not null"`
	Password  string    `gorm:"size:255;not null"`
	Email     string    `gorm:"size:128"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (UserPO) TableName() string {
	return "users"
}