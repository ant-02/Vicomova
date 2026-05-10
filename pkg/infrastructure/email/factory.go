package email

import (
	"fmt"

	"vicomova/pkg/config"
)

type EmailType string

const (
	EmailTypeAliyun EmailType = "aliyun"
)

func NewEmailService(emailType EmailType, cfg interface{}) (EmailService, error) {
	switch emailType {
	case EmailTypeAliyun:
		aliyunCfg, ok := cfg.(*config.AliyunEmail)
		if !ok {
			return nil, fmt.Errorf("invalid aliyun email config")
		}
		return NewAliyunEmailService(aliyunCfg)
	default:
		return nil, fmt.Errorf("unsupported email type: %s", emailType)
	}
}
