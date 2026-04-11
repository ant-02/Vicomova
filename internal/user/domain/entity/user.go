package entity

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type User struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string    `json:"username" gorm:"uniqueIndex;size:64;not null"`
	Password  string    `json:"-" gorm:"size:255;not null"`
	Email     string    `json:"email" gorm:"size:128"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

func NewUser(username, password, email string) *User {
	return &User{
		Username: username,
		Password: hashPassword(password),
		Email:    email,
	}
}

func (u *User) CanLogin(password string) bool {
	return u.Password == hashPassword(password)
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}