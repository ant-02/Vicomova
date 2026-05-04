package errors

var (
	ErrUserNotFound     = NewBizError(404, "User not found")
	ErrUserExists       = NewBizError(409, "User already exists")
	ErrPasswordWrong    = NewBizError(401, "Password is incorrect")
	ErrInvalidUsername  = NewBizError(400, "Invalid username")
	ErrInvalidPassword = NewBizError(400, "Invalid password")
	ErrInvalidEmail    = NewBizError(400, "Invalid email")
)
