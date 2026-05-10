package sms

import "fmt"

type SMSType string

const (
	SMSTypeAliyun SMSType = "aliyun"
)

func NewSMSService(smsType SMSType, cfg interface{}) (SMSService, error) {
	switch smsType {
	case SMSTypeAliyun:
		aliyunCfg, ok := cfg.(*AliyunSMSConfig)
		if !ok {
			return nil, fmt.Errorf("invalid aliyun sms config")
		}
		return NewAliyunSMSService(aliyunCfg)
	default:
		return nil, fmt.Errorf("unsupported sms type: %s", smsType)
	}
}
