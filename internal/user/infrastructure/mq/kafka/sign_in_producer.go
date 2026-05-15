package kafka

import (
	"context"
	"encoding/json"

	pkgKafka "vicomova/pkg/infrastructure/kafka"
	"vicomova/pkg/log"
)

// SignInEvent 用户签到事件
type SignInEvent struct {
	UserID          int64  `json:"user_id"`
	SignDate        string `json:"sign_date"` // 格式: 2006-01-02
	ConsecutiveDays int    `json:"consecutive_days"`
}

type SignInProducer struct {
	sender *pkgKafka.Sender
	topic  string
}

func NewSignInProducer(sender *pkgKafka.Sender, topic string) *SignInProducer {
	return &SignInProducer{
		sender: sender,
		topic:  topic,
	}
}

func (p *SignInProducer) Record(ctx context.Context, event *SignInEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		log.Error.Printf("SignInProducer: marshal failed: %v", err)
		return err
	}

	msg := &pkgKafka.Message{
		Key:   nil,
		Value: data,
	}
	if err := p.sender.Send(ctx, pkgKafka.GetBrokers(), p.topic, msg); err != nil {
		log.Error.Printf("SignInProducer: send failed, userID=%d: %v", event.UserID, err)
		return err
	}
	log.Debug.Printf("SignInProducer: sent sign-in event, userID=%d, topic=%s", event.UserID, p.topic)
	return nil
}
