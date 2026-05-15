package kafka

import (
	"context"
	"encoding/json"
	"log"

	infraKafka "vicomova/internal/video/infrastructure/mq/kafka"
)

type Consumer struct {
	consumer *infraKafka.ConsumerManager
}

func NewConsumer(consumer *infraKafka.ConsumerManager) *Consumer {
	return &Consumer{consumer: consumer}
}

func (c *Consumer) Start(ctx context.Context) error {
	signInHandler := NewSignInHandler()
	c.consumer.Register(signInHandler)

	if err := c.consumer.StartAll(); err != nil {
		log.Printf("Consumer.Start: failed to start: %v", err)
		return err
	}

	log.Printf("Consumer.Start: all handlers started")
	return nil
}

var _ interface {
	Start(ctx context.Context) error
} = (*Consumer)(nil)

var _ = json.Marshal // import check
