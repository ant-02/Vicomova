package entity

import (
	"time"

	"vicomova/internal/user/domain/valueobject"
)

type User struct {
	ID        int64
	Username  *valueobject.Username
	Password  *valueobject.Password
	Email     *valueobject.Email
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(username *valueobject.Username, password *valueobject.Password, email *valueobject.Email) *User {
	return &User{
		Username: username,
		Password: password,
		Email:    email,
	}
}
