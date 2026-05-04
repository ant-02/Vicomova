package valueobject

import (
	"errors"
	"regexp"
)

var ErrInvalidUsername = errors.New("invalid username")

type Username struct {
	value string
}

func NewUsername(value string) (*Username, error) {
	if err := validateUsername(value); err != nil {
		return nil, err
	}
	return &Username{value: value}, nil
}

func validateUsername(value string) error {
	if len(value) < 3 || len(value) > 32 {
		return errors.New("username must be between 3 and 32 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(value) {
		return errors.New("username can only contain letters, numbers, and underscores")
	}
	return nil
}

func (u *Username) Value() string {
	return u.value
}

func (u *Username) String() string {
	return u.value
}

func (u *Username) Equal(other *Username) bool {
	return u.value == other.value
}