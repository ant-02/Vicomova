package email

import "context"

type EmailService interface {
	Send(ctx context.Context, toEmail, subject, body string) error
}
