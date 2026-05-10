package email

import (
	"context"
	"fmt"

	"vicomova/pkg/config"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dm20151123 "github.com/alibabacloud-go/dm-20151123/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
)

type AliyunEmailService struct {
	client      *dm20151123.Client
	accountName string
}

func NewAliyunEmailService(cfg *config.Email) (*AliyunEmailService, error) {
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

func (s *AliyunEmailService) Send(ctx context.Context, toEmail, subject, body string) error {
	request := &dm20151123.SingleSendMailRequest{}
	request.AccountName = tea.String(s.accountName)
	request.AddressType = tea.Int32(1)
	request.ToAddress = tea.String(toEmail)
	request.Subject = tea.String(subject)
	request.HtmlBody = tea.String(body)
	request.ReplyToAddress = tea.Bool(true)

	runtime := &util.RuntimeOptions{}
	_, err := s.client.SingleSendMailWithOptions(request, runtime)
	if err != nil {
		return fmt.Errorf("send mail failed: %w", err)
	}

	return nil
}
