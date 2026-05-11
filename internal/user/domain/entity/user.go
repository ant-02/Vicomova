package entity

import (
	"time"

	"vicomova/internal/user/domain/valueobject"

	"gorm.io/gorm"
)

type User struct {
	ID        int64
	Username  *valueobject.Username
	Password  *valueobject.Password
	Email     *valueobject.Email
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func NewUser(username *valueobject.Username, password *valueobject.Password, email *valueobject.Email) *User {
	return &User{
		Username: username,
		Password: password,
		Email:    email,
	}
}
