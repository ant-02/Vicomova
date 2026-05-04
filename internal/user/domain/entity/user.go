package entity

import (
	"time"

	"vicomova/internal/user/domain/valueobject"
)

type User struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  *valueobject.Username `json:"-" gorm:"-"`
	Password  *valueobject.Password `json:"-" gorm:"-"`
	Email     *valueobject.Email    `json:"-" gorm:"-"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}

func NewUser(username *valueobject.Username, password *valueobject.Password, email *valueobject.Email) *User {
	return &User{
		Username: username,
		Password: password,
		Email:    email,
	}
}

func (u *User) CanLogin(plainPassword string) bool {
	return u.Password.Verify(plainPassword)
}