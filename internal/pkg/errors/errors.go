package errors

import (
	"fmt"
)

type BizError struct {
	Code int32
	Msg  string
}

func (e *BizError) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

func NewBizError(code int32, msg string) *BizError {
	return &BizError{Code: code, Msg: msg}
}

var (
	ErrInvalidParams   = NewBizError(400, "Invalid parameters")
	ErrUnauthorized    = NewBizError(401, "Unauthorized")
	ErrUserNotFound    = NewBizError(404, "User not found")
	ErrUserExists      = NewBizError(409, "User already exists")
	ErrPasswordWrong   = NewBizError(401, "Password is incorrect")
	ErrInternalServer  = NewBizError(500, "Internal server error")
	ErrInvalidToken    = NewBizError(401, "Invalid token")
	ErrTokenExpired    = NewBizError(401, "Token expired")
)
