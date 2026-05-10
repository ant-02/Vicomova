package sms

import (
	"context"
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
)

type AliyunSMSConfig struct {
	AccessKey  string `mapstructure:"access_key"`
	AccessSecret string `mapstructure:"access_secret"`
	SignName  string `mapstructure:"sign_name"`
	TemplateCode string `mapstructure:"template_code"`
	Region    string `mapstructure:"region"`
}

type AliyunSMSService struct {
	client      *dysmsapi.Client
	signName    string
	templateCode string
}

func NewAliyunSMSService(cfg *AliyunSMSConfig) (*AliyunSMSService, error) {
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
		RegionId:   tea.String(cfg.Region),
	}
	cfg2.Endpoint = tea.String("dysmsapi.aliyuncs.com")

	client, err := dysmsapi.NewClient(cfg2)
	if err != nil {
		return nil, fmt.Errorf("create dysmsapi client failed: %w", err)
	}

	return &AliyunSMSService{
		client:      client,
		signName:    cfg.SignName,
		templateCode: cfg.TemplateCode,
	}, nil
}

func (s *AliyunSMSService) SendCode(ctx context.Context, phone, code string) error {
	templateParam := fmt.Sprintf(`{"code":"%s"}`, code)

	request := &dysmsapi.SendSmsRequest{}
	request.PhoneNumbers = tea.String(phone)
	request.SignName = tea.String(s.signName)
	request.TemplateCode = tea.String(s.templateCode)
	request.TemplateParam = tea.String(templateParam)

	runtime := &util.RuntimeOptions{}
	_, err := s.client.SendSmsWithOptions(request, runtime)
	if err != nil {
		return fmt.Errorf("send sms failed: %w", err)
	}

	return nil
}
