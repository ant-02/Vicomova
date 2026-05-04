package valueobject

import (
	"errors"
	"regexp"
)

var ErrInvalidEmail = errors.New("invalid email")

type Email struct {
	value string
}

func NewEmail(value string) (*Email, error) {
	if err := validateEmail(value); err != nil {
		return nil, err
	}
	return &Email{value: value}, nil
}

func validateEmail(value string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		return errors.New("invalid email format")
	}
	return nil
}

func (e *Email) Value() string {
	return e.value
}

func (e *Email) String() string {
	return e.value
}

func (e *Email) Equal(other *Email) bool {
	return e.value == other.value
}