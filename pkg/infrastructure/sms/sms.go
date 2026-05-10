package sms

import "context"

type SMSService interface {
	SendCode(ctx context.Context, phone, code string) error
}
