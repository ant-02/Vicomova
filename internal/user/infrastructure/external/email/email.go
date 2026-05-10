package email

import (
	"context"
	"fmt"
	"time"

	emailPkg "vicomova/pkg/infrastructure/email"

	"vicomova/pkg/constants"
)

type VerificationEmailService struct {
	emailSvc emailPkg.EmailService
}

func NewVerificationEmailService(emailSvc emailPkg.EmailService) *VerificationEmailService {
	return &VerificationEmailService{emailSvc: emailSvc}
}

func (s *VerificationEmailService) SendVerificationCode(ctx context.Context, toEmail, code string) error {
	subject := "Vicomova 邮箱验证码"
	emailBody := fmt.Sprintf(`
		<html>
		<body>
			<p>您正在进行安全验证，本次请求的验证码是：</p>
			<h2 style="color: #ff5722;">%s</h2>
			<p>验证码有效期为%d分钟，请勿泄露给他人。</p>
		</body>
		</html>
	`, code, constants.EmailCodeTTL/time.Minute)

	return s.emailSvc.Send(ctx, toEmail, subject, emailBody)
}
