package email

import (
	"context"
	"fmt"
	"time"

	"vicomova/internal/shared/pkg/constants"
	"vicomova/pkg/config"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dm20151123 "github.com/alibabacloud-go/dm-20151123/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
)

type EmailService interface {
	SendVerificationCode(ctx context.Context, toEmail, code string) error
}

type AliyunEmailService struct {
	client      *dm20151123.Client
	accountName string
}

func NewAliyunEmailService(cfg *config.AliyunEmailConfig) (*AliyunEmailService, error) {
	credConfig := &credential.Config{
		AccessKeyId:     tea.String(cfg.AccessKey),
		AccessKeySecret: tea.String(cfg.AccessSecret),
		Type:            tea.String("access_key"),
	}
	cred, err := credential.NewCredential(credConfig)
	if err != nil {
		return nil, fmt.Errorf("create credential failed: %w", err)
	}

	cfg2 := &openapi.Config{
		Credential: cred,
	}
	cfg2.Endpoint = tea.String("dm.aliyuncs.com")
	cfg2.RegionId = tea.String(cfg.Region)

	client, err := dm20151123.NewClient(cfg2)
	if err != nil {
		return nil, fmt.Errorf("create dm client failed: %w", err)
	}

	return &AliyunEmailService{
		client:      client,
		accountName: cfg.AccountName,
	}, nil
}

func (s *AliyunEmailService) SendVerificationCode(ctx context.Context, toEmail, code string) error {
	emailBody := fmt.Sprintf(`
		<html>
		<body>
			<p>您正在进行安全验证，本次请求的验证码是：</p>
			<h2 style="color: #ff5722;">%s</h2>
			<p>验证码有效期为%d分钟，请勿泄露给他人。</p>
		</body>
		</html>
	`, code, constants.EmailCodeTTL/time.Minute)

	request := &dm20151123.SingleSendMailRequest{}
	request.AccountName = tea.String(s.accountName)
	request.AddressType = tea.Int32(1)
	request.ToAddress = tea.String(toEmail)
	request.Subject = tea.String("Vicomova 邮箱验证码")
	request.HtmlBody = tea.String(emailBody)
	request.ReplyToAddress = tea.Bool(true)

	runtime := &util.RuntimeOptions{}
	_, err := s.client.SingleSendMailWithOptions(request, runtime)
	if err != nil {
		return fmt.Errorf("send mail failed: %w", err)
	}

	return nil
}
