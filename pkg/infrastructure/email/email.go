package email

import "context"

type EmailService interface {
	SendVerificationCode(ctx context.Context, toEmail, code string) error
}
