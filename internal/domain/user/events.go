package user

import "time"

type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
}

type UserRegisteredEvent struct {
	UserID   int64
	Username string
	Email    string
}

func (e *UserRegisteredEvent) EventName() string {
	return "user.registered"
}

func (e *UserRegisteredEvent) OccurredAt() time.Time {
	return time.Now()
}

type UserLoggedInEvent struct {
	UserID   int64
	Username string
}

func (e *UserLoggedInEvent) EventName() string {
	return "user.logged_in"
}

func (e *UserLoggedInEvent) OccurredAt() time.Time {
	return time.Now()
}

type UserLoggedOutEvent struct {
	UserID int64
}

func (e *UserLoggedOutEvent) EventName() string {
	return "user.logged_out"
}

func (e *UserLoggedOutEvent) OccurredAt() time.Time {
	return time.Now()
}

type TokenRefreshedEvent struct {
	UserID        int64
	OldTokenID    string
	NewTokenID    string
}

func (e *TokenRefreshedEvent) EventName() string {
	return "user.token_refreshed"
}

func (e *TokenRefreshedEvent) OccurredAt() time.Time {
	return time.Now()
}
