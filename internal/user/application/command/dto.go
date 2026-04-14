package command

// Command DTOs

type LoginCommand struct {
	Username string
	Password string
}

type RefreshTokenCommand struct {
	RefreshToken string
}

type LogoutCommand struct {
	AccessToken string
}

type SendVerificationCodeCommand struct {
	Username string
	Password string
	Email    string
}

type VerifyAndRegisterCommand struct {
	Username string
	Password string
	Email    string
	Code     string
}

// Result DTOs

type TokenResult struct {
	UserID       int64
	Username     string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

type UserResult struct {
	UserID   int64
	Username string
	Email    string
}

type SendCodeResult struct {
	Success bool
	Message string
}
